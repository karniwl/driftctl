package aws

import (
	"github.com/cloudskiff/driftctl/pkg/remote/aws/repository"
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"
)

type ApiGatewayV2VpcLinkEnumerator struct {
	repository repository.ApiGatewayV2Repository
	factory    resource.ResourceFactory
}

func NewApiGatewayV2VpcLinkEnumerator(repo repository.ApiGatewayV2Repository, factory resource.ResourceFactory) *ApiGatewayV2VpcLinkEnumerator {
	return &ApiGatewayV2VpcLinkEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *ApiGatewayV2VpcLinkEnumerator) SupportedType() resource.ResourceType {
	return aws.AwsApiGatewayV2VpcLinkResourceType
}

func (e *ApiGatewayV2VpcLinkEnumerator) Enumerate() ([]*resource.Resource, error) {
	vpcLinks, err := e.repository.ListAllVpcLinks()
	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(vpcLinks))

	for _, vpcLink := range vpcLinks {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				*vpcLink.VpcLinkId,
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
