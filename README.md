# Mdcodec

Transform Go structures into readable Markdown, tailored for both human and machine consumption.

## Features

- Hierarchical struct representation using `#` for primary and nested lists for secondary structures.
- Uses reflection.

## Roadmap

- Introduce concise single-line encodings.
- Introduce markdown table support (See [moul/mdtable](https://github.com/moul/mdtable)).
- Integrate native support within Amino.
- Examples: develop markdown-centric APIs in Gnoland contracts. And Go clients.

## Examples

### `mdcodec.`Marshal``

```go
type Person struct {
    Name    string
    Age     int
    Address struct {
        City  string
        State string
    }
}

p := Person{}
p.Name = "John Doe"
p.Age = 30
p.Address.City = "Sprintfield"
p.Address.State = "IL"

md := mdcodec.Marshal(p)
fmt.Println(md)

// Output:
// # John Doe (Person)
// - **Age**: 30
// - **Address**:
//   - **City**: Springfield
//   - **State**: IL
```

### `mdcodec.Unmarshal`

```go
var p Person
err := mdcodec.Unmarshal(md, &p)
```

## 🔧 Installation

    go get moul.io/mdcodec
