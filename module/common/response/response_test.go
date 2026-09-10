package response

import "testing"

func TestNormalizeUsesLowerCamelCase(t *testing.T) {
	type record struct {
		ID            string
		PublicID      string
		PublishStatus string
	}
	value := Normalize(record{ID: "1", PublicID: "p", PublishStatus: "draft"})
	result, ok := value.(map[string]interface{})
	if !ok {
		t.Fatalf("Normalize() type = %T", value)
	}
	if result["id"] != "1" || result["publicID"] != "p" || result["publishStatus"] != "draft" {
		t.Fatalf("Normalize() = %#v", result)
	}
}
