package cluster

import (
	"fmt"

	tftypes "github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/protobuf/types/known/structpb"
	"gopkg.in/yaml.v3"
)

// helmValuesForState returns the tftypes.String value to store in state for
// extra_helm_values after a Read. If the existing state value is semantically
// equivalent to the API response, it is preserved unchanged so that the
// original user-provided format (YAML or JSON) is not replaced with JSON.
func helmValuesForState(apiValues *structpb.Struct, stateValue tftypes.String) (tftypes.String, error) {
	// Treat nil/empty API values as "{}" so they can be compared consistently
	// with state values like "" or "{}" that are also semantically empty.
	var apiJSON string
	if apiValues != nil && len(apiValues.Fields) > 0 {
		jsonBytes, err := apiValues.MarshalJSON()
		if err != nil {
			return tftypes.StringNull(), fmt.Errorf("could not marshal extra_helm_values to JSON: %w", err)
		}
		apiJSON = string(jsonBytes)
	} else {
		apiJSON = "{}"
	}

	// If the state already holds a value, keep it when it is semantically
	// equivalent to what the API returned. This avoids replacing YAML with
	// JSON and triggering spurious diffs on the next plan, and prevents
	// empty values like "" or "{}" from being replaced with null.
	if !stateValue.IsNull() && !stateValue.IsUnknown() {
		stateNormalized, err := normalizeHelmValuesToJSON(stateValue.ValueString())
		if err == nil && stateNormalized == apiJSON {
			return stateValue, nil
		}
	}

	// If the API had no values and the state wasn't preserved, return null
	// rather than "{}" so that unset attributes remain unset.
	if apiJSON == "{}" {
		return tftypes.StringNull(), nil
	}

	return tftypes.StringValue(apiJSON), nil
}

// normalizeHelmValuesToJSON parses a YAML or JSON string and re-marshals it via
// structpb so the output matches the JSON form returned by the API.
func normalizeHelmValuesToJSON(input string) (string, error) {
	var helmValues map[string]any
	if err := yaml.Unmarshal([]byte(input), &helmValues); err != nil {
		return "", fmt.Errorf("invalid YAML/JSON: %w", err)
	}

	helmStruct, err := structpb.NewStruct(helmValues)
	if err != nil {
		return "", fmt.Errorf("failed to convert to Struct: %w", err)
	}

	jsonBytes, err := helmStruct.MarshalJSON()
	if err != nil {
		return "", fmt.Errorf("failed to marshal to JSON: %w", err)
	}

	return string(jsonBytes), nil
}
