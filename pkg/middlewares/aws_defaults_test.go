package middlewares

import (
	"strings"
	"testing"

	awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/r3labs/diff/v2"

	"github.com/cloudskiff/driftctl/pkg/resource"
	"github.com/cloudskiff/driftctl/pkg/resource/aws"
)

func TestAwsIamRoleDefaults_Execute(t *testing.T) {
	tests := []struct {
		name               string
		remoteResources    []resource.Resource
		resourcesFromState []resource.Resource
		expected           diff.Changelog
	}{
		{
			"default iam roles when they're not managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "AWSServiceRoleForSSO",
					Path: func(path string) *string { return &path }("/aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(path string) *string { return &path }("/not-aws-service-role/sso.amazonaws.com/"),
				},
				&aws.AwsRoute{
					Id:           "dummy-route",
					RouteTableId: awssdk.String("default-route-table"),
					GatewayId:    awssdk.String("local"),
				},
			},
			[]resource.Resource{
				&aws.AwsRoute{
					Id:           "dummy-route",
					RouteTableId: awssdk.String("default-route-table"),
					GatewayId:    awssdk.String("local"),
				},
			},
			diff.Changelog{
				{
					Type: "delete",
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "OrganizationAccountAccessRole",
						Path: func(path string) *string { return &path }("/not-aws-service-role/sso.amazonaws.com/"),
					},
					To: nil,
				},
			},
		},
		{
			"default iam roles when they're managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:          "AWSServiceRoleForSSO",
					Path:        func(path string) *string { return &path }("/aws-service-role/sso.amazonaws.com/"),
					Description: func(path string) *string { return &path }("test"),
				},
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(path string) *string { return &path }("/not-aws-service-role/sso.amazonaws.com/"),
				},
				&aws.AwsIamRole{
					Id: "driftctl_assume_role:driftctl_policy.10",
					Tags: map[string]string{
						"test": "value",
					},
				},
			},
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "AWSServiceRoleForSSO",
					Path: func(path string) *string { return &path }("/aws-service-role/sso.amazonaws.com/"),
				},
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(path string) *string { return &path }("/not-aws-service-role/sso.amazonaws.com/"),
				},
				&aws.AwsIamRole{
					Id:   "driftctl_assume_role:driftctl_policy.10",
					Tags: map[string]string{},
				},
			},
			diff.Changelog{
				{
					Type: diff.UPDATE,
					Path: []string{"0", "Description"},
					From: func(path string) *string { return &path }("test"),
					To:   nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2", "Tags", "test"},
					From: "value",
					To:   nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := awsIamRoleDefaults(&tt.remoteResources, &tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			changelog, err := diff.Diff(tt.remoteResources, tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			diffs, err := diff.Diff(tt.expected, changelog)
			if err != nil {
				t.Fatal(err)
			}
			if len(diffs) == 0 {
				return
			}

			for _, change := range diffs {
				t.Errorf("%s got = %v, want %v", strings.Join(change.Path, "."), awsutil.Prettify(change.From), awsutil.Prettify(change.To))
			}
		})
	}
}

func TestAwsIamPolicyAttachmentDefaults_Execute(t *testing.T) {
	tests := []struct {
		name               string
		remoteResources    []resource.Resource
		resourcesFromState []resource.Resource
		expected           diff.Changelog
	}{
		{
			"test that default iam policy attachment are excluded when not managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "custom-role",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRole{
					Id:   "AWSServiceRoleForSSO",
					Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
					Roles: &[]string{"custom-role"},
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/AWSSSOServiceRolePolicy",
					Roles: &[]string{"AWSServiceRoleForSSO"},
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/whatever",
					Roles: &[]string{"AWSServiceRoleForSSO"},
				},
			},
			[]resource.Resource{},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "custom-role",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamRole{
						Id:   "AWSServiceRoleForSSO",
						Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamPolicyAttachment{
						Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
						Roles: &[]string{"custom-role"},
					},
					To: nil,
				},
			},
		},
		{
			"test that default iam policy attachment are excluded when not managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "custom-role",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRole{
					Id:   "AWSServiceRoleForSSO",
					Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
					Roles: &[]string{"custom-role"},
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/AWSSSOServiceRolePolicy",
					Roles: &[]string{"AWSServiceRoleForSSO"},
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/whatever",
					Roles: &[]string{"AWSServiceRoleForSSO", "custom-role"},
				},
			},
			[]resource.Resource{
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/AWSSSOServiceRolePolicy",
					Roles: &[]string{"AWSServiceRoleForSSO"},
				},
			},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "custom-role",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamRole{
						Id:   "AWSServiceRoleForSSO",
						Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamPolicyAttachment{
						Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
						Roles: &[]string{"custom-role"},
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"4"},
					From: &aws.AwsIamPolicyAttachment{
						Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/whatever",
						Roles: &[]string{"AWSServiceRoleForSSO", "custom-role"},
					},
					To: nil,
				},
			},
		},
		{
			"do not ignore default iam policy attachment when role cannot be found",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "custom-role",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
					Roles: &[]string{"custom-role"},
				},
				&aws.AwsIamPolicyAttachment{
					Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/AWSSSOServiceRolePolicy",
					Roles: &[]string{"AWSServiceRoleForSSO"},
				},
			},
			[]resource.Resource{},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "custom-role",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamPolicyAttachment{
						Id:    "driftctl_test-arn:aws:iam::0123456789:policy/driftctl",
						Roles: &[]string{"custom-role"},
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamPolicyAttachment{
						Id:    "AWSServiceRoleForSSO-arn:aws:iam::aws:policy/aws-service-role/AWSSSOServiceRolePolicy",
						Roles: &[]string{"AWSServiceRoleForSSO"},
					},
					To: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := awsIamPolicyAttachmentDefaults(&tt.remoteResources, &tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			changelog, err := diff.Diff(tt.remoteResources, tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			diffs, err := diff.Diff(tt.expected, changelog)
			if err != nil {
				t.Fatal(err)
			}
			if len(diffs) == 0 {
				return
			}

			for _, change := range diffs {
				t.Errorf("%s got = %v, want %v", strings.Join(change.Path, "."), awsutil.Prettify(change.From), awsutil.Prettify(change.To))
			}
		})
	}
}

func TestAwsIamRolePolicyDefaults_Execute(t *testing.T) {
	tests := []struct {
		name               string
		remoteResources    []resource.Resource
		resourcesFromState []resource.Resource
		expected           diff.Changelog
	}{
		{
			"ignore default iam role policies when they're managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "AWSServiceRoleForSSO",
					Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "AWSServiceRoleForSSO:AdministratorAccess",
					Role: func(p string) *string { return &p }("AWSServiceRoleForSSO"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "OrganizationAccountAccessRole:AdministratorAccess",
					Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
				},
				&aws.AwsRoute{
					Id:           "dummy-route",
					RouteTableId: awssdk.String("default-route-table"),
					GatewayId:    awssdk.String("local"),
				},
			},
			[]resource.Resource{
				&aws.AwsRoute{
					Id:           "dummy-route",
					RouteTableId: awssdk.String("default-route-table"),
					GatewayId:    awssdk.String("local"),
				},
			},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "AWSServiceRoleForSSO",
						Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamRole{
						Id:   "OrganizationAccountAccessRole",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamRolePolicy{
						Id:   "OrganizationAccountAccessRole:AdministratorAccess",
						Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
					},
					To: nil,
				},
			},
		},
		{
			"iam role policies when they're managed by IaC",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "custom-role",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "driftctl_assume_role:driftctl_policy.10",
					Role: func(p string) *string { return &p }("custom-role"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "OrganizationAccountAccessRole:AdministratorAccess",
					Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
				},
			},
			[]resource.Resource{
				&aws.AwsIamRolePolicy{
					Id:   "OrganizationAccountAccessRole:AdministratorAccess",
					Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
				},
			},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "custom-role",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamRole{
						Id:   "OrganizationAccountAccessRole",
						Path: func(p string) *string { return &p }("/aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamRolePolicy{
						Id:   "driftctl_assume_role:driftctl_policy.10",
						Role: func(p string) *string { return &p }("custom-role"),
					},
					To: nil,
				},
			},
		},
		{
			"do not ignore default iam role policies when role cannot be found",
			[]resource.Resource{
				&aws.AwsIamRole{
					Id:   "OrganizationAccountAccessRole",
					Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "AWSServiceRoleForSSO:AdministratorAccess",
					Role: func(p string) *string { return &p }("AWSServiceRoleForSSO"),
				},
				&aws.AwsIamRolePolicy{
					Id:   "OrganizationAccountAccessRole:AdministratorAccess",
					Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
				},
			},
			[]resource.Resource{},
			diff.Changelog{
				{
					Type: diff.DELETE,
					Path: []string{"0"},
					From: &aws.AwsIamRole{
						Id:   "OrganizationAccountAccessRole",
						Path: func(p string) *string { return &p }("/not-aws-service-role/sso.amazonaws.com"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"1"},
					From: &aws.AwsIamRolePolicy{
						Id:   "AWSServiceRoleForSSO:AdministratorAccess",
						Role: func(p string) *string { return &p }("AWSServiceRoleForSSO"),
					},
					To: nil,
				},
				{
					Type: diff.DELETE,
					Path: []string{"2"},
					From: &aws.AwsIamRolePolicy{
						Id:   "OrganizationAccountAccessRole:AdministratorAccess",
						Role: func(p string) *string { return &p }("OrganizationAccountAccessRole"),
					},
					To: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := awsIamRolePolicyDefaults(&tt.remoteResources, &tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			changelog, err := diff.Diff(tt.remoteResources, tt.resourcesFromState)
			if err != nil {
				t.Fatal(err)
			}

			diffs, err := diff.Diff(tt.expected, changelog)
			if err != nil {
				t.Fatal(err)
			}
			if len(diffs) == 0 {
				return
			}

			for _, change := range diffs {
				t.Errorf("%s got = %v, want %v", strings.Join(change.Path, "."), awsutil.Prettify(change.From), awsutil.Prettify(change.To))
			}
		})
	}
}