package exchangepolicy

import (
	"context"
	"fmt"
	"strings"

	exchangepolicypb "github.com/cofide/cofide-api-sdk/gen/go/proto/exchange_policy/v1alpha1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
)

const actionPrefix = "EXCHANGE_POLICY_ACTION_"

var stringMatcherAttrTypes = map[string]attr.Type{
	"exact": tftypes.StringType,
	"glob":  tftypes.StringType,
}

var stringMatcherObjectType = tftypes.ObjectType{
	AttrTypes: stringMatcherAttrTypes,
}

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

	return proto, nil
}

// protoToModel converts an ExchangePolicy protobuf to an equivalent ExchangePolicyModel.
func protoToModel(proto *exchangepolicypb.ExchangePolicy) ExchangePolicyModel {
	model := ExchangePolicyModel{
		ID:              tftypes.StringValue(proto.GetId()),
		OrgID:           tftypes.StringValue(proto.GetOrgId()),
		Name:            tftypes.StringValue(proto.GetName()),
		TrustZoneID:     tftypes.StringValue(proto.GetTrustZoneId()),
		SubjectIdentity: stringSetFromProto(proto.GetSubjectIdentity()),
		SubjectIssuer:   stringSetFromProto(proto.GetSubjectIssuer()),
		ActorIdentity:   stringSetFromProto(proto.GetActorIdentity()),
		ActorIssuer:     stringSetFromProto(proto.GetActorIssuer()),
		ClientID:        stringSetFromProto(proto.GetClientId()),
		TargetAudience:  stringSetFromProto(proto.GetTargetAudience()),
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

	return model
}

// stringSetToProto converts a tftypes.List of StringMatcherModel to a StringSet protobuf.
func stringSetToProto(ctx context.Context, list tftypes.List) (*exchangepolicypb.StringSet, error) {
	if list.IsNull() || list.IsUnknown() {
		return nil, nil
	}
	var matchers []StringMatcherModel
	list.ElementsAs(ctx, &matchers, false)
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
// Returns an error if neither Exact nor Glob is set.
func stringMatcherToProto(model StringMatcherModel) (*exchangepolicypb.StringMatcher, error) {
	if !model.Exact.IsNull() {
		return &exchangepolicypb.StringMatcher{
			Match: &exchangepolicypb.StringMatcher_Exact{
				Exact: model.Exact.ValueString(),
			},
		}, nil
	}
	if !model.Glob.IsNull() {
		return &exchangepolicypb.StringMatcher{
			Match: &exchangepolicypb.StringMatcher_Glob{
				Glob: model.Glob.ValueString(),
			},
		}, nil
	}
	return nil, fmt.Errorf("string matcher must set exactly one of exact or glob")
}

// stringSetFromProto converts a StringSet protobuf to a tftypes.List of StringMatcherModel.
// An empty StringSet (no matchers) is treated as absent and returns a null list.
func stringSetFromProto(proto *exchangepolicypb.StringSet) tftypes.List {
	if proto == nil || len(proto.GetMatchers()) == 0 {
		return tftypes.ListNull(stringMatcherObjectType)
	}
	elems := make([]attr.Value, 0, len(proto.GetMatchers()))
	for _, m := range proto.GetMatchers() {
		mm := stringMatcherFromProto(m)
		elems = append(elems, tftypes.ObjectValueMust(stringMatcherAttrTypes, map[string]attr.Value{
			"exact": mm.Exact,
			"glob":  mm.Glob,
		}))
	}
	return tftypes.ListValueMust(stringMatcherObjectType, elems)
}

// stringMatcherFromProto converts a StringMatcher protobuf to a StringMatcherModel.
func stringMatcherFromProto(proto *exchangepolicypb.StringMatcher) StringMatcherModel {
	model := StringMatcherModel{
		Exact: tftypes.StringNull(),
		Glob:  tftypes.StringNull(),
	}
	switch m := proto.GetMatch().(type) {
	case *exchangepolicypb.StringMatcher_Exact:
		model.Exact = tftypes.StringValue(m.Exact)
	case *exchangepolicypb.StringMatcher_Glob:
		model.Glob = tftypes.StringValue(m.Glob)
	}
	return model
}
