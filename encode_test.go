package mdcodec

import "testing"

type Person struct {
	Name    string
	Age     int
	Address struct {
		Street string
		City   string
	}
}

func TestMarshal(t *testing.T) {
	p := Person{Name: "John Doe", Age: 30}
	p.Address.Street = "123 Maple St."
	p.Address.City = "Springfield"

	md := Marshal(p)
	expected := `# Person

- **Name**: John Doe
- **Age**: 30

## Address

- **Street**: 123 Maple St.
- **City**: Springfield

`

	if md != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, md)
	}
}

func TestUnmarshal(t *testing.T) {
	md := `# Person

- **Name**: John Doe
- **Age**: 30

## Address

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
