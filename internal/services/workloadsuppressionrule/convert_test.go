package workloadsuppressionrule

import (
	"context"
	"testing"
	"time"

	workloadsuppressionrulepb "github.com/cofide/cofide-api-sdk/gen/go/proto/workload_suppression_rule/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestProtoToModel_Minimal(t *testing.T) {
	proto := &workloadsuppressionrulepb.WorkloadSuppressionRule{
		Id:    "rule-1",
		OrgId: "org-1",
		Name:  "test-rule",
	}

	got := protoToModel(proto)

	assert.Equal(t, types.StringValue("rule-1"), got.ID)
	assert.Equal(t, types.StringValue("org-1"), got.OrgID)
	assert.Equal(t, types.StringValue("test-rule"), got.Name)
	assert.True(t, got.Description.IsNull())
	assert.Equal(t, types.BoolValue(false), got.Enabled)
	assert.Nil(t, got.KubernetesPod)
	assert.True(t, got.CreatedAt.IsNull())
	assert.True(t, got.LastUpdatedAt.IsNull())
}

func TestProtoToModel_Full(t *testing.T) {
	createdAt := time.Date(2026, 7, 29, 12, 0, 0, 0, time.UTC)
	lastUpdatedAt := time.Date(2026, 7, 30, 8, 30, 0, 0, time.UTC)

	proto := &workloadsuppressionrulepb.WorkloadSuppressionRule{
		Id:          "rule-2",
		OrgId:       "org-2",
		Name:        "full-rule",
		Description: "suppress noisy debug pods",
		Enabled:     true,
		Matcher: &workloadsuppressionrulepb.WorkloadSuppressionRule_KubernetesPod{
			KubernetesPod: &workloadsuppressionrulepb.KubernetesPodMatcher{
				TrustZoneIds: []string{"tz-1"},
				ClusterIds:   []string{"cluster-1", "cluster-2"},
				Namespaces:   []string{"kube-system"},
				Labels:       map[string]string{"app": "debug"},
			},
		},
		CreatedAt:     timestamppb.New(createdAt),
		LastUpdatedAt: timestamppb.New(lastUpdatedAt),
	}

	got := protoToModel(proto)

	assert.Equal(t, types.StringValue("rule-2"), got.ID)
	assert.Equal(t, types.StringValue("org-2"), got.OrgID)
	assert.Equal(t, types.StringValue("full-rule"), got.Name)
	assert.Equal(t, types.StringValue("suppress noisy debug pods"), got.Description)
	assert.Equal(t, types.BoolValue(true), got.Enabled)
	require.NotNil(t, got.KubernetesPod)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tz-1")}), got.KubernetesPod.TrustZoneIDs)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("cluster-1"), types.StringValue("cluster-2")}), got.KubernetesPod.ClusterIDs)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("kube-system")}), got.KubernetesPod.Namespaces)
	assert.Equal(t, types.MapValueMust(types.StringType, map[string]attr.Value{"app": types.StringValue("debug")}), got.KubernetesPod.Labels)
	assert.Equal(t, types.StringValue(createdAt.Format(time.RFC3339)), got.CreatedAt)
	assert.Equal(t, types.StringValue(lastUpdatedAt.Format(time.RFC3339)), got.LastUpdatedAt)
}

func TestModelToProto_KubernetesPodMatcher(t *testing.T) {
	ctx := context.Background()

	model := WorkloadSuppressionRuleModel{
		ID:          types.StringValue("rule-3"),
		OrgID:       types.StringValue("org-3"),
		Name:        types.StringValue("kube-rule"),
		Description: types.StringValue("desc"),
		Enabled:     types.BoolValue(true),
		KubernetesPod: &KubernetesPodMatcherModel{
			TrustZoneIDs: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("tz-1")}),
			ClusterIDs:   types.ListNull(types.StringType),
			Namespaces:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("default")}),
			Labels:       types.MapValueMust(types.StringType, map[string]attr.Value{"app": types.StringValue("debug")}),
		},
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError(), diags)

	assert.Equal(t, "rule-3", got.GetId())
	assert.Equal(t, "org-3", got.GetOrgId())
	assert.Equal(t, "kube-rule", got.GetName())
	assert.Equal(t, "desc", got.GetDescription())
	assert.True(t, got.GetEnabled())

	pod := got.GetKubernetesPod()
	require.NotNil(t, pod)
	assert.Equal(t, []string{"tz-1"}, pod.GetTrustZoneIds())
	assert.Empty(t, pod.GetClusterIds())
	assert.Equal(t, []string{"default"}, pod.GetNamespaces())
	assert.Equal(t, map[string]string{"app": "debug"}, pod.GetLabels())
}

func TestModelToProto_NoMatcher(t *testing.T) {
	ctx := context.Background()

	model := WorkloadSuppressionRuleModel{
		ID:    types.StringValue("rule-4"),
		OrgID: types.StringValue("org-4"),
		Name:  types.StringValue("no-matcher-rule"),
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError(), diags)
	assert.Nil(t, got.GetMatcher())
}

func TestProtoToModel_AWSLambdaFunctionMatcher(t *testing.T) {
	proto := &workloadsuppressionrulepb.WorkloadSuppressionRule{
		Id:    "rule-5",
		OrgId: "org-5",
		Name:  "lambda-rule",
		Matcher: &workloadsuppressionrulepb.WorkloadSuppressionRule_AwsLambdaFunction{
			AwsLambdaFunction: &workloadsuppressionrulepb.AWSLambdaFunctionMatcher{
				CloudAccountIds: []string{"acct-1"},
				Regions:         []string{"us-east-1", "us-west-2"},
				FunctionNames:   []string{"my-function"},
				Tags:            map[string]string{"env": "prod"},
			},
		},
	}

	got := protoToModel(proto)

	require.NotNil(t, got.AWSLambdaFunction)
	assert.Nil(t, got.KubernetesPod)
	assert.Nil(t, got.AWSAgentCoreRuntime)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("acct-1")}), got.AWSLambdaFunction.CloudAccountIDs)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("us-east-1"), types.StringValue("us-west-2")}), got.AWSLambdaFunction.Regions)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-function")}), got.AWSLambdaFunction.FunctionNames)
	assert.Equal(t, types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}), got.AWSLambdaFunction.Tags)
}

func TestModelToProto_AWSLambdaFunctionMatcher(t *testing.T) {
	ctx := context.Background()

	model := WorkloadSuppressionRuleModel{
		ID:    types.StringValue("rule-6"),
		OrgID: types.StringValue("org-6"),
		Name:  types.StringValue("lambda-rule"),
		AWSLambdaFunction: &AWSLambdaFunctionMatcherModel{
			CloudAccountIDs: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("acct-1")}),
			Regions:         types.ListNull(types.StringType),
			FunctionNames:   types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-function")}),
			Tags:            types.MapValueMust(types.StringType, map[string]attr.Value{"env": types.StringValue("prod")}),
		},
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError(), diags)

	lambda := got.GetAwsLambdaFunction()
	require.NotNil(t, lambda)
	assert.Equal(t, []string{"acct-1"}, lambda.GetCloudAccountIds())
	assert.Empty(t, lambda.GetRegions())
	assert.Equal(t, []string{"my-function"}, lambda.GetFunctionNames())
	assert.Equal(t, map[string]string{"env": "prod"}, lambda.GetTags())
}

func TestProtoToModel_AWSAgentCoreRuntimeMatcher(t *testing.T) {
	proto := &workloadsuppressionrulepb.WorkloadSuppressionRule{
		Id:    "rule-7",
		OrgId: "org-7",
		Name:  "agentcore-rule",
		Matcher: &workloadsuppressionrulepb.WorkloadSuppressionRule_AwsAgentcoreRuntime{
			AwsAgentcoreRuntime: &workloadsuppressionrulepb.AWSAgentCoreRuntimeMatcher{
				CloudAccountIds:   []string{"acct-2"},
				Regions:           []string{"eu-west-1"},
				AgentRuntimeNames: []string{"my-agent", "other-agent"},
			},
		},
	}

	got := protoToModel(proto)

	require.NotNil(t, got.AWSAgentCoreRuntime)
	assert.Nil(t, got.KubernetesPod)
	assert.Nil(t, got.AWSLambdaFunction)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("acct-2")}), got.AWSAgentCoreRuntime.CloudAccountIDs)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("eu-west-1")}), got.AWSAgentCoreRuntime.Regions)
	assert.Equal(t, types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-agent"), types.StringValue("other-agent")}), got.AWSAgentCoreRuntime.AgentRuntimeNames)
}

func TestModelToProto_AWSAgentCoreRuntimeMatcher(t *testing.T) {
	ctx := context.Background()

	model := WorkloadSuppressionRuleModel{
		ID:    types.StringValue("rule-8"),
		OrgID: types.StringValue("org-8"),
		Name:  types.StringValue("agentcore-rule"),
		AWSAgentCoreRuntime: &AWSAgentCoreRuntimeMatcherModel{
			CloudAccountIDs:   types.ListNull(types.StringType),
			Regions:           types.ListValueMust(types.StringType, []attr.Value{types.StringValue("eu-west-1")}),
			AgentRuntimeNames: types.ListValueMust(types.StringType, []attr.Value{types.StringValue("my-agent")}),
		},
	}

	got, diags := modelToProto(ctx, model)
	require.False(t, diags.HasError(), diags)

	runtime := got.GetAwsAgentcoreRuntime()
	require.NotNil(t, runtime)
	assert.Empty(t, runtime.GetCloudAccountIds())
	assert.Equal(t, []string{"eu-west-1"}, runtime.GetRegions())
	assert.Equal(t, []string{"my-agent"}, runtime.GetAgentRuntimeNames())
}
