// Package mdcodec provides tools to convert Go structures to and from readable Markdown.
//
// The primary functionality of the package is to offer marshalling and unmarshalling capabilities
// between Go structures and a corresponding Markdown representation. This Markdown representation
// can be both machine-readable and visually appealing for human readers.
//
// Example:
//
//	type Person struct {
//	    Name    string
//	    Age     int
//	    Address struct {
//		City  string
//		State string
//	    }
//	}
//
//	p := Person{
//	    Name: "John Doe",
//	    Age:  30,
//	    Address: struct {
//		City  string
//		State string
//	    }{
//		City:  "Springfield",
//		State: "IL",
//	    },
//	}
//
//	md := Marshal(p)
//	// This will produce:
//	//
//	// # Person
//	//
//	// - **Name**: John Doe
//	// - **Age**: 30
//	//
//	// ## Address
//	//
//	// - **City**: Springfield
//	// - **State**: IL
package mdcodec

import "fmt"

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

	md := Marshal(p)
	fmt.Println(md)
	// Output:
	// # Person
	//
	// - **Name**: John Doe
	// - **Age**: 30
	//
	// ## Address
	//
	// - **City**: Springfield
	// - **State**: IL
}

func ExampleUnmarshal() {
	type Person struct {
		Name    string
		Age     int
		Address struct {
			City  string
			State string
		}
	}

	md := `
# Person

- **Name**: John Doe
- **Age**: 30

## Address

- **City**: Springfield
- **State**: IL
`

	var p Person
	Unmarshal(md, &p)
	fmt.Printf("Name: %s, Age: %d, City: %s, State: %s", p.Name, p.Age, p.Address.City, p.Address.State)
	// Output:
	// Name: John Doe, Age: 30, City: Springfield, State: IL
}
