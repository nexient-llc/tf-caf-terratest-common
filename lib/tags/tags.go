package tags

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func ExtendExpectedTagsByThoseAddedByFramework(expectedTags map[string]string, resourceName string) map[string]string {
	expectedTags["resource_name"] = resourceName
	return expectedTags
}

func CompareExpectedTagsVsActual(t *testing.T, expectedTags map[string]string, actualTags map[string]string) {
	for k, v := range expectedTags {
		assert.Equal(t, v, actualTags[k]) //dont (falsely) alarm when a 3rd party tool added a tag to the cloud resource at runtime
	}
}
