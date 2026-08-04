package trustzoneserver

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewUpdateMaskCoversAllFields guards against the update mask silently
// drifting from the generated proto type: if a new field is added to
// UpdateTrustZoneServerRequest_UpdateMask (e.g. because a new resource
// attribute was added) but newUpdateMask isn't updated to set it, Connect
// will never receive changes to that field on update.
func TestNewUpdateMaskCoversAllFields(t *testing.T) {
	mask := newUpdateMask()

	v := reflect.ValueOf(mask).Elem()
	typ := v.Type()

	var unset []string
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() || field.Type.Kind() != reflect.Bool {
			continue
		}
		if !v.Field(i).Bool() {
			unset = append(unset, field.Name)
		}
	}

	assert.Empty(t, unset, "newUpdateMask must set every mask field to true; found unset field(s) %v — a new field was likely added to the proto UpdateMask without being wired up here", unset)
}
