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
//	md, _ := Marshal(p)
//	// This will produce:
//	//
//	// # John Doe (Person)
//	//
//	// - **Name**: John Doe
//	// - **Age**: 30
//	//
//	// ## Address
//	//
//	// - **City**: Springfield
//	// - **State**: IL
package mdcodec
