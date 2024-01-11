package utils

type Options struct {
    Host string
    Port string
}

// NewOptions returns a new initialized Options object
func NewOptions() *Options {
	return &Options{}
}
