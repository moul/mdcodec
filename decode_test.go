package mdcodec

import (
	"fmt"
	"testing"
)

func ExampleUnmarshal() {
	type Person struct {
		Name    string `md:"title"`
		Age     int
		Address struct {
			City  string
			State string
		}
	}

	md := `# John Doe (Person)

- **Age**: 30
- **Address**:
  - **City**: Springfield
  - **State**: IL
`

	var p Person
	Unmarshal(md, &p)
	fmt.Printf("Name: %s, Age: %d, City: %s, State: %s", p.Name, p.Age, p.Address.City, p.Address.State)
	// Output:
	// Name: John Doe, Age: 30, City: Springfield, State: IL
}

func TestUnmarshal(t *testing.T) {
	md := `# John Doe (Person)

- **Age**: 30
- **Address**:
  - **Street**: 123 Maple St.
  - **City**: Springfield
`

	var p Person
	err := Unmarshal(md, &p)
	if err != nil {
		t.Fatalf("Error during unmarshal: %v", err)
	}

	if p.Name != "John Doe" {
		t.Errorf("Expected Name 'John Doe', Got %s", p.Name)
	}

	if p.Age != 30 {
		t.Errorf("Expected Age 30, Got %d", p.Age)
	}

	if p.Address.Street != "123 Maple St." {
		t.Errorf("Expected Street '123 Maple St.', Got %s", p.Address.Street)
	}

	if p.Address.City != "Springfield" {
		t.Errorf("Expected City 'Springfield', Got %s", p.Address.City)
	}
}

func TestInvalidUnmarshal(t *testing.T) {
	md := `# John Doe (Person)

- **Age**: Thirty
`

	var p Person
	err := Unmarshal(md, &p)
	if err == nil {
		t.Error("Expected error due to invalid age format, got nil")
	}
}
