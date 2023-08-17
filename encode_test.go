package mdcodec

import (
	"fmt"
	"testing"
)

func ExampleMarshal() {
	type Person struct {
		Name    string
		Age     int
		Address struct {
			City  string
			State string
		}
	}

	p := Person{
		Name: "John Doe",
		Age:  30,
		Address: struct {
			City  string
			State string
		}{
			City:  "Springfield",
			State: "IL",
		},
	}

	md, _ := Marshal(p)
	fmt.Println(md)
	// Output:
	// # John Doe (Person)
	// - **Age**: 30
	// - **Address**:
	//   - **City**: Springfield
	//   - **State**: IL
}

func TestMarshal(t *testing.T) {
	p := Person{Name: "John Doe", Age: 30}
	p.Address.Street = "123 Maple St."
	p.Address.City = "Springfield"

	md, _ := Marshal(p)
	expected := `# John Doe (Person)
- **Age**: 30
- **Address**:
  - **Street**: 123 Maple St.
  - **City**: Springfield
`

	if md != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, md)
	}
}
