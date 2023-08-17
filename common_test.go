package mdcodec

type Person struct {
	Name    string
	Age     int
	Address struct {
		Street string
		City   string
	}
}
