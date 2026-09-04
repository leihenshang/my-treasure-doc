package response

import "testing"

func TestNormalizeUsesLowerCamelCase(t *testing.T) {
	type record struct {
		ID            string
		PublishStatus string
	}
	value := normalize(record{ID: "1", PublishStatus: "draft"})
	result, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("normalize() type = %T", value)
	}
	if result["id"] != "1" || result["publishStatus"] != "draft" {
		t.Fatalf("normalize() = %#v", result)
	}
}
