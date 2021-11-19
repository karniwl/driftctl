package repository

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewayv2"
	"github.com/aws/aws-sdk-go/service/apigatewayv2/apigatewayv2iface"
	"github.com/cloudskiff/driftctl/pkg/remote/cache"
)

type ApiGatewayV2Repository interface {
	ListAllVpcLinks() ([]*apigatewayv2.VpcLink, error)
}

type apigatewayv2Repository struct {
	client apigatewayv2iface.ApiGatewayV2API
	cache  cache.Cache
}

func NewApiGatewayV2Repository(session *session.Session, c cache.Cache) *apigatewayv2Repository {
	return &apigatewayv2Repository{
		apigatewayv2.New(session),
		c,
	}
}

func (r *apigatewayv2Repository) ListAllVpcLinks() ([]*apigatewayv2.VpcLink, error) {
	if v := r.cache.Get("apigatewayv2ListAllVpcLinks"); v != nil {
		return v.([]*apigatewayv2.VpcLink), nil
	}

	input := apigatewayv2.GetVpcLinksInput{}
	resources, err := r.client.GetVpcLinks(&input)
	if err != nil {
		return nil, err
	}

	r.cache.Put("apigatewayv2ListAllVpcLinks", resources.Items)
	return resources.Items, nil
}
