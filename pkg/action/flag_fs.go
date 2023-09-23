//go:build ignore

package action

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"syscall"

	"github.com/urfave/cli/v2"
)

type FSFlag struct {
	Name string

	Category    string
	DefaultText string
	FilePath    string
	Usage       string

	Required   bool
	Hidden     bool
	HasBeenSet bool

	Value       fs.FS
	Destination *fs.FS

	Aliases []string
	EnvVars []string

	defaultValue fs.FS

	TakesFile bool

	Action func(*cli.Context, fs.FS) error
}

var (
	_ cli.Flag = (*FSFlag)(nil)
)

// String returns a readable representation of this value (for usage defaults)
func (f *FSFlag) String() string {
	return cli.FlagStringer(f)
}

// IsSet returns whether or not the flag has been set through env or file
func (f *FSFlag) IsSet() bool {
	return f.HasBeenSet
}

// Names returns the names of the flag
func (f *FSFlag) Names() []string {
	return cli.FlagNames(f.Name, f.Aliases)
}

// IsRequired returns whether or not the flag is required
func (f *FSFlag) IsRequired() bool {
	return f.Required
}

// IsVisible returns true if the flag is not hidden, otherwise false
func (f *FSFlag) IsVisible() bool {
	return !f.Hidden
}

// TakesValue returns true of the flag takes a value, otherwise false
func (f *FSFlag) TakesValue() bool {
	return true
}

// GetUsage returns the usage string for the flag
func (f *FSFlag) GetUsage() string {
	return f.Usage
}

// GetCategory returns the category for the flag
func (f *FSFlag) GetCategory() string {
	return f.Category
}

// GetEnvVars returns the env vars for this flag
func (f *FSFlag) GetEnvVars() []string {
	return f.EnvVars
}

// Apply populates the flag given the flag set and environment
func (f *FSFlag) Apply(set *flag.FlagSet) error {
	// set default value so that environment wont be able to overwrite it
	f.defaultValue = f.Value

	if val, _, found := flagFromEnvOrFile(f.EnvVars, f.FilePath); found {
		f.Value = val
		f.HasBeenSet = true
	}

	for _, name := range f.Names() {
		if f.Destination != nil {
			set.StringVar(f.Destination, name, f.Value, f.Usage)
			continue
		}
		set.String(name, f.Value, f.Usage)
	}

	return nil
}

// Get returns the flag’s value in the given Context.
func (f *FSFlag) Get(ctx *cli.Context) string {
	return ctx.Path(f.Name)
}

// RunAction executes flag action if set
func (f *FSFlag) RunAction(c *cli.Context) error {
	if f.Action != nil {
		return f.Action(c, c.Path(f.Name))
	}

	return nil
}

func flagFromEnvOrFile(envVars []string, filePath string) (value string, fromWhere string, found bool) {
	for _, envVar := range envVars {
		envVar = strings.TrimSpace(envVar)
		if value, found := syscall.Getenv(envVar); found {
			return value, fmt.Sprintf("environment variable %q", envVar), true
		}
	}
	for _, fileVar := range strings.Split(filePath, ",") {
		if fileVar != "" {
			if data, err := os.ReadFile(fileVar); err == nil {
				return string(data), fmt.Sprintf("file %q", filePath), true
			}
		}
	}
	return "", "", false
}
