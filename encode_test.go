package mdcodec

import (
	"fmt"
	"testing"
)

func ExampleMarshal() {
	type Person struct {
		Name    string `md:"title"`
		Age     int
		Address struct {
			Street string
			City   string
		}
	}

	p := Person{}
	p.Name = "John Doe"
	p.Age = 30
	p.Address.Street = "123 Maple St."
	p.Address.City = "Springfield"

	md, _ := Marshal(p)
	fmt.Println(md)
	// Output:
	// # John Doe (Person)
	//
	// - **Age**: 30
	// - **Address**:
	//   - **Street**: 123 Maple St.
	//   - **City**: Springfield
}

func TestMarshal(t *testing.T) {
	p := Person{}
	p.Name = "John Doe"
	p.Age = 30
	p.Address.Street = "123 Maple St."
	p.Address.City = "Springfield"

	md, err := Marshal(p)
	if err != nil {
		t.Fatalf("Error during Marshal: %v", err)
	}

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
