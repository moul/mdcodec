package mdcodec

type Person struct {
	Name    string `md:"title"`
	Age     int
	Address struct {
		Street string
		City   string
	}
}
