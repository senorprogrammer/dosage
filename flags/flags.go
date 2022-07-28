package flags

import (
	"fmt"

	goFlags "github.com/jessevdk/go-flags"
)

// Flags is a struct that contains all the flags passed in via the command line on startup
type Flags struct {
	APIKey string `short:"k" long:"apikey" required:"yes" description:"Your DigitalOcean API key"`
}

// NewFlags creates an instance of Flags
func NewFlags() *Flags {
	flags := Flags{}
	return &flags
}

// Parse parses the incoming flags
func (flags *Flags) Parse() error {
	parser := goFlags.NewParser(flags, goFlags.Default)
	if _, err := parser.Parse(); err != nil {
		return err
	}

	fmt.Println(flags.APIKey)

	return nil
}
