package google

import (
	remoteerror "github.com/cloudskiff/driftctl/pkg/remote/error"
	"github.com/cloudskiff/driftctl/pkg/remote/google/repository"
	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/google"
)

type GoogleComputeInstanceEnumerator struct {
	repository repository.AssetRepository
	factory    resource.ResourceFactory
}

func NewGoogleComputeInstanceEnumerator(repo repository.AssetRepository, factory resource.ResourceFactory) *GoogleComputeInstanceEnumerator {
	return &GoogleComputeInstanceEnumerator{
		repository: repo,
		factory:    factory,
	}
}

func (e *GoogleComputeInstanceEnumerator) SupportedType() resource.ResourceType {
	return google.GoogleComputeInstanceResourceType
}

func (e *GoogleComputeInstanceEnumerator) Enumerate() ([]*resource.Resource, error) {
	resources, err := e.repository.SearchAllInstances()

	if err != nil {
		return nil, remoteerror.NewResourceListingError(err, string(e.SupportedType()))
	}

	results := make([]*resource.Resource, 0, len(resources))

	for _, res := range resources {
		results = append(
			results,
			e.factory.CreateAbstractResource(
				string(e.SupportedType()),
				trimResourceName(res.GetName()),
				map[string]interface{}{},
			),
		)
	}

	return results, err
}
