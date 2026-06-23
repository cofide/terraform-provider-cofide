package exchangepolicy

import (
	"context"
	"fmt"
	"strings"
	"time"

	exchangepolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/exchange_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/durationpb"
)

const actionPrefix = "EXCHANGE_POLICY_ACTION_"

var stringMatcherAttrTypes = map[string]attr.Type{
	"exact": tftypes.StringType,
	"glob":  tftypes.StringType,
}

var stringMatcherObjectType = tftypes.ObjectType{
	AttrTypes: stringMatcherAttrTypes,
}

// spiffeMtlsAttrTypes holds the attribute types for the spiffe_mtls auth variant.
var spiffeMtlsAttrTypes = map[string]attr.Type{
	"spiffe_id": tftypes.StringType,
}

// authAttrTypes holds one entry per supported auth variant. When a new proto
// oneof variant is added, add its attr type here and handle it in the convert
// functions below; ExternalHookModel itself does not need to change.
var authAttrTypes = map[string]attr.Type{
	"spiffe_mtls": tftypes.ObjectType{AttrTypes: spiffeMtlsAttrTypes},
}

var externalHookAttrTypes = map[string]attr.Type{
	"name":        tftypes.StringType,
	"description": tftypes.StringType,
	"url":         tftypes.StringType,
	"auth":        tftypes.ObjectType{AttrTypes: authAttrTypes},
	"timeout":     tftypes.Int64Type,
}

var externalHookObjectType = tftypes.ObjectType{AttrTypes: externalHookAttrTypes}

// modelToProto converts an ExchangePolicyModel to an equivalent ExchangePolicy protobuf.
func modelToProto(ctx context.Context, model ExchangePolicyModel) (*exchangepolicypb.ExchangePolicy, error) {
	proto := &exchangepolicypb.ExchangePolicy{
		Id:          model.ID.ValueString(),
		Name:        model.Name.ValueString(),
		TrustZoneId: model.TrustZoneID.ValueString(),
	}

	if !model.Action.IsNull() && !model.Action.IsUnknown() {
		key := actionPrefix + model.Action.ValueString()
		val, ok := exchangepolicypb.ExchangePolicyAction_value[key]
		if !ok {
			return nil, fmt.Errorf("invalid action %q: must be one of ALLOW, DENY", model.Action.ValueString())
		}
		action := exchangepolicypb.ExchangePolicyAction(val)
		proto.Action = &action
	}

	fields := []struct {
		name  string
		list  tftypes.List
		dest  **exchangepolicypb.StringSet
	}{
		{"subject_identity", model.SubjectIdentity, &proto.SubjectIdentity},
		{"subject_issuer", model.SubjectIssuer, &proto.SubjectIssuer},
		{"actor_identity", model.ActorIdentity, &proto.ActorIdentity},
		{"actor_issuer", model.ActorIssuer, &proto.ActorIssuer},
		{"subject_audience", model.SubjectAudience, &proto.SubjectAudience},
		{"client_id", model.ClientID, &proto.ClientId},
		{"target_audience", model.TargetAudience, &proto.TargetAudience},
	}
	for _, f := range fields {
		ss, err := stringSetToProto(ctx, f.list)
		if err != nil {
			return nil, fmt.Errorf("field %s: %w", f.name, err)
		}
		*f.dest = ss
	}

	if !model.OutboundScopes.IsNull() && !model.OutboundScopes.IsUnknown() {
		for _, v := range model.OutboundScopes.Elements() {
			if sv, ok := v.(tftypes.String); ok {
				proto.OutboundScopes = append(proto.OutboundScopes, sv.ValueString())
			}
		}
	}

	hooks, err := externalHooksToProto(ctx, model.ExternalHooks)
	if err != nil {
		return nil, err
	}
	proto.ExternalHooks = hooks

	return proto, nil
}

// protoToModel converts an ExchangePolicy protobuf to an equivalent ExchangePolicyModel.
func protoToModel(proto *exchangepolicypb.ExchangePolicy) (ExchangePolicyModel, error) {
	type stringSetField struct {
		name string
		ss   *exchangepolicypb.StringSet
		dest *tftypes.List
	}

	model := ExchangePolicyModel{
		ID:          tftypes.StringValue(proto.GetId()),
		OrgID:       tftypes.StringValue(proto.GetOrgId()),
		Name:        tftypes.StringValue(proto.GetName()),
		TrustZoneID: tftypes.StringValue(proto.GetTrustZoneId()),
	}

	fields := []stringSetField{
		{"subject_identity", proto.GetSubjectIdentity(), &model.SubjectIdentity},
		{"subject_issuer", proto.GetSubjectIssuer(), &model.SubjectIssuer},
		{"actor_identity", proto.GetActorIdentity(), &model.ActorIdentity},
		{"actor_issuer", proto.GetActorIssuer(), &model.ActorIssuer},
		{"subject_audience", proto.GetSubjectAudience(), &model.SubjectAudience},
		{"client_id", proto.GetClientId(), &model.ClientID},
		{"target_audience", proto.GetTargetAudience(), &model.TargetAudience},
	}
	for _, f := range fields {
		list, err := stringSetFromProto(f.ss)
		if err != nil {
			return ExchangePolicyModel{}, fmt.Errorf("field %s: %w", f.name, err)
		}
		*f.dest = list
	}

	if proto.Action != nil {
		model.Action = tftypes.StringValue(strings.TrimPrefix(proto.GetAction().String(), actionPrefix))
	} else {
		model.Action = tftypes.StringNull()
	}

	scopes := make([]attr.Value, len(proto.GetOutboundScopes()))
	for i, scope := range proto.GetOutboundScopes() {
		scopes[i] = tftypes.StringValue(scope)
	}
	model.OutboundScopes = tftypes.ListValueMust(tftypes.StringType, scopes)

	model.ExternalHooks = externalHooksFromProto(proto.GetExternalHooks())

	return model, nil
}

// stringSetToProto converts a tftypes.List of StringMatcherModel to a StringSet protobuf.
func stringSetToProto(ctx context.Context, list tftypes.List) (*exchangepolicypb.StringSet, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var matchers []StringMatcherModel
	if diags := list.ElementsAs(ctx, &matchers, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert string matchers: %v", diags)
	}
	proto := &exchangepolicypb.StringSet{}
	for _, m := range matchers {
		matcher, err := stringMatcherToProto(m)
		if err != nil {
			return nil, err
		}
		proto.Matchers = append(proto.Matchers, matcher)
	}
	return proto, nil
}

// stringMatcherToProto converts a StringMatcherModel to a StringMatcher protobuf.
// Returns an error if exactly one of Exact or Glob is not set to a known value.
func stringMatcherToProto(model StringMatcherModel) (*exchangepolicypb.StringMatcher, error) {
	exactSet := !model.Exact.IsNull() && !model.Exact.IsUnknown()
	globSet := !model.Glob.IsNull() && !model.Glob.IsUnknown()
	if exactSet == globSet {
		return nil, fmt.Errorf("string matcher must set exactly one of exact or glob")
	}
	if exactSet {
		return &exchangepolicypb.StringMatcher{
			Match: &exchangepolicypb.StringMatcher_Exact{
				Exact: model.Exact.ValueString(),
			},
		}, nil
	}
	return &exchangepolicypb.StringMatcher{
		Match: &exchangepolicypb.StringMatcher_Glob{
			Glob: model.Glob.ValueString(),
		},
	}, nil
}

// stringSetFromProto converts a StringSet protobuf to a tftypes.List of StringMatcherModel.
// An empty StringSet (no matchers) is treated as absent and returns a null list.
func stringSetFromProto(proto *exchangepolicypb.StringSet) (tftypes.List, error) {
	if proto == nil || len(proto.GetMatchers()) == 0 {
		return tftypes.ListNull(stringMatcherObjectType), nil
	}
	elems := make([]attr.Value, 0, len(proto.GetMatchers()))
	for _, m := range proto.GetMatchers() {
		mm, err := stringMatcherFromProto(m)
		if err != nil {
			return tftypes.List{}, err
		}
		elems = append(elems, tftypes.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
			"exact": mm.Exact,
			"glob":  mm.Glob,
		}))
	}
	return tftypes.ListValueMust(stringMatcherObjectType, elems), nil
}

// stringMatcherFromProto converts a StringMatcher protobuf to a StringMatcherModel.
// Returns an error if neither Exact nor Glob is set.
func stringMatcherFromProto(proto *exchangepolicypb.StringMatcher) (StringMatcherModel, error) {
	switch m := proto.GetMatch().(type) {
	case *exchangepolicypb.StringMatcher_Exact:
		return StringMatcherModel{
			Exact: tftypes.StringValue(m.Exact),
			Glob:  tftypes.StringNull(),
		}, nil
	case *exchangepolicypb.StringMatcher_Glob:
		return StringMatcherModel{
			Exact: tftypes.StringNull(),
			Glob:  tftypes.StringValue(m.Glob),
		}, nil
	default:
		return StringMatcherModel{}, fmt.Errorf("string matcher must set exactly one of exact or glob")
	}
}

// externalHooksToProto converts the ExternalHooks list in a model to a slice
// of ExternalHook protobufs. Returns nil when the list is null or unknown.
func externalHooksToProto(ctx context.Context, list tftypes.List) ([]*exchangepolicypb.ExternalHook, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var hooks []ExternalHookModel
	if diags := list.ElementsAs(ctx, &hooks, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert external hooks: %v", diags)
	}
	result := make([]*exchangepolicypb.ExternalHook, 0, len(hooks))
	for i, h := range hooks {
		pb, err := externalHookToProto(h)
		if err != nil {
			return nil, fmt.Errorf("external_hook[%d]: %w", i, err)
		}
		result = append(result, pb)
	}
	return result, nil
}

// externalHookToProto converts a single ExternalHookModel to an ExternalHook protobuf.
func externalHookToProto(model ExternalHookModel) (*exchangepolicypb.ExternalHook, error) {
	pb := &exchangepolicypb.ExternalHook{
		Name: model.Name.ValueString(),
		Url:  model.URL.ValueString(),
	}
	if !model.Description.IsNull() && !model.Description.IsUnknown() {
		pb.Description = model.Description.ValueString()
	}
	if !model.Auth.IsNull() && !model.Auth.IsUnknown() {
		// Each auth variant is a separate attribute inside the auth object.
		// When a new proto oneof variant is added, add a new case here.
		authAttrs := model.Auth.Attributes()
		if sm, ok := authAttrs["spiffe_mtls"].(tftypes.Object); ok && !sm.IsNull() && !sm.IsUnknown() {
			if spiffeIDAttr, ok := sm.Attributes()["spiffe_id"].(tftypes.String); ok && !spiffeIDAttr.IsNull() && !spiffeIDAttr.IsUnknown() {
				pb.Auth = &exchangepolicypb.ExternalHook_SpiffeMtls{
					SpiffeMtls: &exchangepolicypb.SpiffeMtlsAuth{SpiffeId: spiffeIDAttr.ValueString()},
				}
			}
		}
	}
	if !model.Timeout.IsNull() && !model.Timeout.IsUnknown() {
		pb.Timeout = durationpb.New(time.Duration(model.Timeout.ValueInt64()) * time.Second)
	}
	return pb, nil
}

// externalHooksFromProto converts a slice of ExternalHook protobufs to a
// tftypes.List of external hook objects. Returns a null list when hooks is empty.
func externalHooksFromProto(hooks []*exchangepolicypb.ExternalHook) tftypes.List {
	if len(hooks) == 0 {
		return tftypes.ListNull(externalHookObjectType)
	}
	elems := make([]attr.Value, 0, len(hooks))
	for _, h := range hooks {
		elems = append(elems, externalHookFromProto(h))
	}
	return tftypes.ListValueMust(externalHookObjectType, elems)
}

// externalHookFromProto converts a single ExternalHook protobuf to a tftypes.Object.
func externalHookFromProto(pb *exchangepolicypb.ExternalHook) tftypes.Object {
	description := tftypes.StringNull()
	if pb.GetDescription() != "" {
		description = tftypes.StringValue(pb.GetDescription())
	}

	// Build the auth object. All auth variant attributes must be present; unused
	// variants are set to null. When a new proto oneof variant is added, add its
	// null value here and populate it in the relevant case below.
	authAttrs := map[string]attr.Value{
		"spiffe_mtls": tftypes.ObjectNull(spiffeMtlsAttrTypes),
	}
	switch a := pb.GetAuth().(type) {
	case *exchangepolicypb.ExternalHook_SpiffeMtls:
		authAttrs["spiffe_mtls"] = tftypes.ObjectValueMust(spiffeMtlsAttrTypes, map[string]attr.Value{
			"spiffe_id": tftypes.StringValue(a.SpiffeMtls.GetSpiffeId()),
		})
	}
	authObj := tftypes.ObjectNull(authAttrTypes)
	if pb.GetAuth() != nil {
		authObj = tftypes.ObjectValueMust(authAttrTypes, authAttrs)
	}

	timeout := tftypes.Int64Null()
	if pb.GetTimeout() != nil {
		timeout = tftypes.Int64Value(int64(pb.GetTimeout().AsDuration().Seconds()))
	}

	return tftypes.ObjectValueMust(externalHookAttrTypes, map[string]attr.Value{
		"name":        tftypes.StringValue(pb.GetName()),
		"description": description,
		"url":         tftypes.StringValue(pb.GetUrl()),
		"auth":        authObj,
		"timeout":     timeout,
	})
}
