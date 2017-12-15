package filter

// Plugin modifies the input, which is a map representation of the json received
// to a output map or errors out.
type Plugin interface {
	Name() string
	Filter(input *Stream) error
}
