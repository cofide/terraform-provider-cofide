package cloudorganization

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFullUpdateMask guards against the update path regressing to a nil
// update_mask, which the API rejects (see fullUpdateMask's doc comment):
// every mutable field must be present and set to true.
func TestFullUpdateMask(t *testing.T) {
	mask := fullUpdateMask()

	assert.True(t, mask.GetName())
	assert.True(t, mask.GetAwsAudience())
	assert.True(t, mask.GetDiscoveryEnabled())
	assert.True(t, mask.GetAwsRoleChain())
	assert.True(t, mask.GetDiscoveryInterval())
	assert.True(t, mask.GetAwsAssumeThroughOidc())
}
