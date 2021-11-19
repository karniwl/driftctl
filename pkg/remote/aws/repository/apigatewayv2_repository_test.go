package repository

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
	awstest "github.com/cloudskiff/driftctl/test/aws"
	"github.com/pkg/errors"

	"github.com/r3labs/diff/v2"
	"github.com/stretchr/testify/assert"
)

func Test_apigatewayv2Repository_ListAllVpcLinks(t *testing.T) {
	vpcLinks := []*apigatewayv2.VpcLink{
		{VpcLinkId: aws.String("vpcLink1")},
		{VpcLinkId: aws.String("vpcLink2")},
		{VpcLinkId: aws.String("vpcLink3")},
		{VpcLinkId: aws.String("vpcLink4")},
		{VpcLinkId: aws.String("vpcLink5")},
		{VpcLinkId: aws.String("vpcLink6")},
	}

	remoteError := errors.New("remote error")

	tests := []struct {
		name    string
		mocks   func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache)
		want    []*apigatewayv2.VpcLink
		wantErr error
	}{
		{
			name: "list multiple vpc links",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetVpcLinks",
					&apigatewayv2.GetVpcLinksInput{}).Return(&apigatewayv2.GetVpcLinksOutput{Items: vpcLinks}, nil).Once()

				store.On("Get", "apigatewayv2ListAllVpcLinks").Return(nil).Times(1)
				store.On("Put", "apigatewayv2ListAllVpcLinks", vpcLinks).Return(false).Times(1)
			},
			want: vpcLinks,
		},
		{
			name: "should hit cache",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				store.On("Get", "apigatewayv2ListAllVpcLinks").Return(vpcLinks).Times(1)
			},
			want: vpcLinks,
		},
		{
			name: "should return remote error",
			mocks: func(client *awstest.MockFakeApiGatewayV2, store *cache.MockCache) {
				client.On("GetVpcLinks",
					&apigatewayv2.GetVpcLinksInput{}).Return(nil, remoteError).Once()

				store.On("Get", "apigatewayv2ListAllVpcLinks").Return(nil).Times(1)
			},
			wantErr: remoteError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store := &cache.MockCache{}
			client := &awstest.MockFakeApiGatewayV2{}
			tt.mocks(client, store)
			r := &apigatewayv2Repository{
				client: client,
				cache:  store,
			}
			got, err := r.ListAllVpcLinks()
			assert.Equal(t, tt.wantErr, err)

			changelog, err := diff.Diff(got, tt.want)
			assert.Nil(t, err)
			if len(changelog) > 0 {
				for _, change := range changelog {
					t.Errorf("%s: %s -> %s", strings.Join(change.Path, "."), change.From, change.To)
				}
				t.Fail()
			}
			store.AssertExpectations(t)
			client.AssertExpectations(t)
		})
	}
}
