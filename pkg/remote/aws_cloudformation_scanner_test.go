package remote

import (
	"errors"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/cloudskiff/driftctl/mocks"
	"github.com/cloudskiff/driftctl/pkg/filter"
	"github.com/cloudskiff/driftctl/pkg/remote/alerts"
	"github.com/cloudskiff/driftctl/pkg/remote/aws"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	remoteerr "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	"github.com/cloudskiff/driftctl/test"
	"github.com/cloudskiff/driftctl/test/goldenfile"
	testresource "github.com/cloudskiff/driftctl/test/resource"
	terraform2 "github.com/cloudskiff/driftctl/test/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCloudformationStack(t *testing.T) {
	tests := []struct {
		test    string
		dirName string
		mocks   func(*repository.MockCloudformationRepository, *mocks.AlerterInterface)
		wantErr error
	}{
		{
			test:    "no cloudformation stacks",
			dirName: "aws_cloudformation_stack_empty",
			mocks: func(repository *repository.MockCloudformationRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllStacks").Return([]*cloudformation.Stack{}, nil)
			},
		},
		{
			test:    "multiple cloudformation stacks",
			dirName: "aws_cloudformation_stack_multiple",
			mocks: func(repository *repository.MockCloudformationRepository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllStacks").Return([]*cloudformation.Stack{
					{StackId: awssdk.String("arn:aws:cloudformation:us-east-1:047081014315:stack/bar-stack/c7a96e70-0f21-11ec-bd2a-0a2d95c2b2ab")},
					{StackId: awssdk.String("arn:aws:cloudformation:us-east-1:047081014315:stack/foo-stack/c7aa0ab0-0f21-11ec-ba25-129d8c0b3757")},
				}, nil)
			},
		},
		{
			test:    "cannot list cloudformation stacks",
			dirName: "aws_cloudformation_stack_list",
			mocks: func(repository *repository.MockCloudformationRepository, alerter *mocks.AlerterInterface) {
				awsError := awserr.NewRequestFailure(awserr.New("AccessDeniedException", "", errors.New("")), 400, "")
				repository.On("ListAllStacks").Return(nil, awsError)

				alerter.On("SendAlert", resourceaws.AwsCloudformationStackResourceType, alerts.NewRemoteAccessDeniedAlert(common.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(awsError, resourceaws.AwsCloudformationStackResourceType, resourceaws.AwsCloudformationStackResourceType), alerts.EnumerationPhase)).Return()
			},
			wantErr: nil,
		},
	}

	schemaRepository := testresource.InitFakeSchemaRepository("aws", "3.19.0")
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)
	deserializer := resource.NewDeserializer(factory)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			shouldUpdate := c.dirName == *goldenfile.Update

			sess := session.Must(session.NewSessionWithOptions(session.Options{
				SharedConfigState: session.SharedConfigEnable,
			}))

			scanOptions := ScannerOptions{Deep: true}
			providerLibrary := terraform.NewProviderLibrary()
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository.MockCloudformationRepository{}
			c.mocks(fakeRepo, alerter)

			var repo repository.CloudformationRepository = fakeRepo
			providerVersion := "3.19.0"
			realProvider, err := terraform2.InitTestAwsProvider(providerLibrary, providerVersion)
			if err != nil {
				t.Fatal(err)
			}
			provider := terraform2.NewFakeTerraformProvider(realProvider)
			provider.WithResponse(c.dirName)

			// Replace mock by real resources if we are in update mode
			if shouldUpdate {
				err := realProvider.Init()
				if err != nil {
					t.Fatal(err)
				}
				provider.ShouldUpdate()
				repo = repository.NewCloudformationRepository(sess, cache.New(0))
			}

			remoteLibrary.AddEnumerator(aws.NewCloudformationStackEnumerator(repo, factory))
			remoteLibrary.AddDetailsFetcher(resourceaws.AwsCloudformationStackResourceType, common.NewGenericDetailsFetcher(resourceaws.AwsCloudformationStackResourceType, provider, deserializer))

			testFilter := &filter.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.wantErr)
			if err != nil {
				return
			}
			test.TestAgainstGoldenFile(got, resourceaws.AwsCloudformationStackResourceType, c.dirName, provider, deserializer, shouldUpdate, tt)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
		})
	}
}
