package mdcodec

import "testing"

func TestInvalidUnmarshal(t *testing.T) {
	md := `# Person

- **Name**: John Doe
- **Age**: Thirty
`

	var p Person
	err := Unmarshal(md, &p)
	if err == nil {
		t.Error("Expected error due to invalid age format, got nil")
	}
}
