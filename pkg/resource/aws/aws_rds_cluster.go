package aws

import (
	"github.com/cloudskiff/driftctl/pkg/resource"
)

const AwsRDSClusterResourceType = "aws_rds_cluster"

func initAwsRDSClusterMetaData(resourceSchemaRepository resource.SchemaRepositoryInterface) {
	resourceSchemaRepository.SetResolveReadAttributesFunc(AwsRDSClusterResourceType, func(res *resource.Resource) map[string]string {
		return map[string]string{
			"cluster_identifier": *res.Attributes().GetString("cluster_identifier"),
			"database_name":      *res.Attributes().GetString("database_name"),
		}
	})
	resourceSchemaRepository.SetNormalizeFunc(AwsRDSClusterResourceType, func(res *resource.Resource) {
		val := res.Attributes()
		val.SafeDelete([]string{"timeouts"})
		val.SafeDelete([]string{"master_password"})
		val.SafeDelete([]string{"cluster_members"})
		val.SafeDelete([]string{"skip_final_snapshot"})
		val.SafeDelete([]string{"allow_major_version_upgrade"})
		val.SafeDelete([]string{"apply_immediately"})
		val.SafeDelete([]string{"final_snapshot_identifier"})
		val.SafeDelete([]string{"source_region"})
	})
	resourceSchemaRepository.SetFlags(AwsRDSClusterResourceType, resource.FlagDeepMode)
}
