package remote

import (
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/cloudskiff/driftctl/mocks"
	"github.com/cloudskiff/driftctl/pkg/filter"
	"github.com/cloudskiff/driftctl/pkg/remote/alerts"
	"github.com/cloudskiff/driftctl/pkg/remote/aws"
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	"github.com/cloudskiff/driftctl/pkg/remote/common"
	remoteerr "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	resourceaws "github.com/cloudskiff/driftctl/pkg/resource/aws"
	"github.com/cloudskiff/driftctl/pkg/terraform"
	testresource "github.com/cloudskiff/driftctl/test/resource"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestApiGatewayV2VpcLink(t *testing.T) {
	dummyError := errors.New("this is an error")

	tests := []struct {
		test           string
		mocks          func(*repository.MockApiGatewayV2Repository, *mocks.AlerterInterface)
		assertExpected func(t *testing.T, got []*resource.Resource)
		wantErr        error
	}{
		{
			test: "no api gateway v2 vpc links",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllVpcLinks").Return([]*apigatewayv2.VpcLink{}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 0)
			},
		},
		{
			test: "single api gateway v2 vpc link",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllVpcLinks").Return([]*apigatewayv2.VpcLink{
					{VpcLinkId: awssdk.String("b8r351")},
				}, nil)
			},
			assertExpected: func(t *testing.T, got []*resource.Resource) {
				assert.Len(t, got, 1)

				assert.Equal(t, got[0].ResourceId(), "b8r351")
				assert.Equal(t, got[0].ResourceType(), resourceaws.AwsApiGatewayV2VpcLinkResourceType)
			},
		},
		{
			test: "cannot list api gateway v2 vpc links",
			mocks: func(repository *repository.MockApiGatewayV2Repository, alerter *mocks.AlerterInterface) {
				repository.On("ListAllVpcLinks").Return(nil, dummyError)
				alerter.On("SendAlert", resourceaws.AwsApiGatewayV2VpcLinkResourceType, alerts.NewRemoteAccessDeniedAlert(common.RemoteAWSTerraform, remoteerr.NewResourceListingErrorWithType(dummyError, resourceaws.AwsApiGatewayV2VpcLinkResourceType, resourceaws.AwsApiGatewayV2VpcLinkResourceType), alerts.EnumerationPhase)).Return()
			},
			wantErr: remoteerr.NewResourceListingError(dummyError, resourceaws.AwsApiGatewayV2VpcLinkResourceType),
		},
	}

	providerVersion := "3.19.0"
	schemaRepository := testresource.InitFakeSchemaRepository("aws", providerVersion)
	resourceaws.InitResourcesMetadata(schemaRepository)
	factory := terraform.NewTerraformResourceFactory(schemaRepository)

	for _, c := range tests {
		t.Run(c.test, func(tt *testing.T) {
			scanOptions := ScannerOptions{}
			remoteLibrary := common.NewRemoteLibrary()

			// Initialize mocks
			alerter := &mocks.AlerterInterface{}
			fakeRepo := &repository.MockApiGatewayV2Repository{}
			c.mocks(fakeRepo, alerter)

			var repo repository.ApiGatewayV2Repository = fakeRepo

			remoteLibrary.AddEnumerator(aws.NewApiGatewayV2VpcLinkEnumerator(repo, factory))

			testFilter := &filter.MockFilter{}
			testFilter.On("IsTypeIgnored", mock.Anything).Return(false)

			s := NewScanner(remoteLibrary, alerter, scanOptions, testFilter)
			got, err := s.Resources()
			assert.Equal(tt, err, c.wantErr)
			if err != nil {
				return
			}

			c.assertExpected(tt, got)
			alerter.AssertExpectations(tt)
			fakeRepo.AssertExpectations(tt)
			testFilter.AssertExpectations(tt)
		})
	}
}
