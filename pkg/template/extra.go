package template

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/helmwave/helmwave/pkg/parallel"
	"gopkg.in/yaml.v3"
)

// Values is alias for string map of interfaces.
type Values = map[string]any

// ToYaml renders data into YAML string.
// Used as custom template function.
func ToYaml(v any) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("failed to marshal %v to YAML: %w", v, err)
	}

	return string(data), nil
}

// FromYaml parses YAML string into data.
// Used as custom template function.
func FromYaml(str string) (Values, error) {
	m := Values{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s from YAML: %w", str, err)
	}

	return m, nil
}

// Exec runs external binary and returns its standard output.
// Used as custom template function.
func Exec(command string, args []any, inputs ...string) (string, error) {
	var input string
	if len(inputs) > 0 {
		input = inputs[0]
	}

	strArgs := make([]string, len(args))
	for i, a := range args {
		switch a := a.(type) {
		case string:
			strArgs[i] = a
		default:
			return "", fmt.Errorf("unexpected type of arg \"%s\" in args %v at index %d", reflect.TypeOf(a), args, i)
		}
	}
	cmd := exec.Command(command, strArgs...)
	// cmd.Dir = c.basePath

	wg := parallel.NewWaitGroup()

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("failed to create stdin pipe for command: %w", err)
	}

	wg.Add(1)
	go writeCommandInput(stdin, input, wg)

	output := &bytes.Buffer{}
	wg.Add(1)
	go getCommandOutput(cmd, output, wg)

	if err := wg.Wait(); err != nil {
		return "", fmt.Errorf("failed to run command: %w", err)
	}

	return output.String(), nil
}

func writeCommandInput(stdin io.WriteCloser, input string, wg *parallel.WaitGroup) {
	defer func(stdin io.WriteCloser, wg *parallel.WaitGroup) {
		wg.ErrChan() <- stdin.Close()
		wg.Done()
	}(stdin, wg)

	i := 0
	size := len(input)

	for i < size {
		n, err := io.WriteString(stdin, input[i:])
		if err != nil {
			wg.ErrChan() <- fmt.Errorf("failed while writing %d bytes to stdin: %w", size, err)

			return
		}

		i += n
	}
}

func getCommandOutput(cmd *exec.Cmd, output *bytes.Buffer, wg *parallel.WaitGroup) {
	defer wg.Done()

	bs, err := cmd.Output()
	if err != nil {
		wg.ErrChan() <- fmt.Errorf("failed to get command output: %w", err)

		return
	}

	_, err = output.Write(bs)
	if err != nil {
		wg.ErrChan() <- fmt.Errorf("failed while copying %d bytes from stdout: %w", len(bs), err)
	}
}

// SetValueAtPath sets value in map by dot-separated key path.
// Used as custom template function.
func SetValueAtPath(path string, value any, values Values) (Values, error) {
	var current any
	current = values
	components := strings.Split(path, ".")
	pathToMap := components[:len(components)-1]
	key := components[len(components)-1]

	for _, k := range pathToMap {
		switch typedCurrent := current.(type) {
		case map[string]any:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path %q: value for key %q does not exist", path, k)
			}
			current = v
		case map[any]any:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path %q: value for key %q does not exist", path, k)
			}
			current = v
		default:
			return nil, fmt.Errorf(
				"failed to walk over path %q: value for key %q is not a map: %v",
				path,
				k,
				reflect.TypeOf(current),
			)
		}
	}

	switch typedCurrent := current.(type) {
	case map[string]any:
		typedCurrent[key] = value
	case map[any]any:
		typedCurrent[key] = value
	default:
		return nil, fmt.Errorf(
			"failed to set value at path %q: value for key %q is not a map: %v",
			path,
			key,
			reflect.TypeOf(current),
		)
	}

	return values, nil
}

// RequiredEnv returns environment variable by name and errors if it is not defined.
// Used as custom template function.
func RequiredEnv(name string) (string, error) {
	if val, exists := os.LookupEnv(name); exists && val != "" {
		return val, nil
	}

	return "", fmt.Errorf("required env var %q is not set", name)
}

// Required returns error if val is nil of empty string. Otherwise it returns the same val.
// Used as custom template function.
func Required(warn string, val any) (any, error) {
	if val == nil {
		return nil, errors.New(warn)
	} else if _, ok := val.(string); ok {
		if val == "" {
			return nil, errors.New(warn)
		}
	}

	return val, nil
}

// ReadFile reads file and returns its contents as string.
// Used as custom template function.
func ReadFile(file string) (string, error) {
	b, err := os.ReadFile(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file %q: %w", file, err)
	}

	return string(b), nil
}

func noKeyError(key string, obj any) error {
	return fmt.Errorf("key %q is not present in %v", key, obj)
}

// Get returns value in map by dot-separated key path.
// First argument is dot-separated key path.
// Second argument is default value if key not found and is optional.
// Third argument is map to search in.
// Used as custom template function.
//
//nolint:gocognit
func Get(path string, varArgs ...any) (any, error) {
	defSet, def, obj, err := parseGetVarArgs(varArgs)
	if err != nil {
		return nil, err
	}

	if path == "" {
		return obj, nil
	}
	keys := strings.Split(path, ".")
	key := keys[0]
	var v any
	var ok bool
	switch typedObj := obj.(type) {
	case map[string]any:
		v, ok = typedObj[key]
		if !ok {
			if defSet {
				return def, nil
			}

			return nil, noKeyError(key, obj)
		}
	case map[any]any:
		v, ok = typedObj[key]
		if !ok {
			if defSet {
				return def, nil
			}

			return nil, noKeyError(key, obj)
		}
	default:
		r, err := tryReflectGet(obj, key, defSet, def)
		if err != nil {
			return nil, err
		}
		v = r
	}

	if defSet {
		return Get(strings.Join(keys[1:], "."), def, v)
	}

	return Get(strings.Join(keys[1:], "."), v)
}

func tryReflectGet(obj any, key string, defSet bool, def any) (any, error) {
	maybeStruct := reflect.ValueOf(obj)
	if maybeStruct.Kind() != reflect.Struct {
		return nil, fmt.Errorf(
			"unexpected type(%v) of value for key %q: it must be either map[string]any or any struct",
			reflect.TypeOf(obj),
			key,
		)
	} else if maybeStruct.NumField() < 1 {
		return nil, noKeyError(key, obj)
	}
	f := maybeStruct.FieldByName(key)
	if !f.IsValid() {
		if defSet {
			return def, nil
		}

		return nil, noKeyError(key, obj)
	}

	return f.Interface(), nil
}

// HasKey searches for any value by dot-separated key path in map.
// Used as custom template function.
func HasKey(path string, varArgs ...any) (bool, error) {
	defSet, def, obj, err := parseGetVarArgs(varArgs)
	if err != nil {
		return false, err
	}

	if path == "" {
		return true, nil
	}
	keys := strings.Split(path, ".")
	var v any
	var ok bool
	switch typedObj := obj.(type) {
	case map[string]any:
		v, ok = typedObj[keys[0]]
		if !ok {
			return defSet, nil
		}
	case map[any]any:
		v, ok = typedObj[keys[0]]
		if !ok {
			return defSet, nil
		}
	default:
		found, f, err := tryReflectHasKey(obj, keys[0], defSet, def)
		if err != nil {
			return false, err
		}
		if f == nil {
			return found, nil
		}

		v = f
	}

	if defSet {
		return HasKey(strings.Join(keys[1:], "."), def, v)
	}

	return HasKey(strings.Join(keys[1:], "."), v)
}

func tryReflectHasKey(obj any, key string, defSet bool, def any) (bool, any, error) {
	maybeStruct := reflect.ValueOf(obj)
	if maybeStruct.Kind() != reflect.Struct {
		return false, nil, fmt.Errorf(
			"unexpected type(%v) of value for key %q: it must be either map[string]any or any struct",
			reflect.TypeOf(obj),
			key,
		)
	} else if maybeStruct.NumField() < 1 {
		return false, nil, noKeyError(key, obj)
	}
	f := maybeStruct.FieldByName(key)
	if !f.IsValid() {
		if defSet {
			return true, def, nil
		}

		return false, nil, noKeyError(key, obj)
	}

	return true, f.Interface(), nil
}

func parseGetVarArgs(varArgs []any) (defSet bool, def, obj any, err error) {
	switch len(varArgs) {
	case 1:
		defSet = false
		def = nil
		obj = varArgs[0]
	case 2:
		defSet = true
		def = varArgs[0]
		obj = varArgs[1]
	default:
		err = fmt.Errorf(
			"unexpected number of args passed to the template function (path, [def, ]obj): "+
				"expected 1 or 2, got %d, args was %q",
			len(varArgs),
			varArgs,
		)
	}

	return
}
