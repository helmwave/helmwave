package template

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/shlex"
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

// FromYamlArray parses YAML array string into data.
// Used as custom template function.
func FromYamlArray(str string) ([]any, error) {
	a := []any{}

	if err := yaml.Unmarshal([]byte(str), &a); err != nil {
		return nil, fmt.Errorf("failed to unmarshal %s from YAML as array: %w", str, err)
	}

	return a, nil
}

// FromYamlAll parses multiple YAML documents separated by `---`.
// Returns an array of all documents.
// Used as custom template function.
func FromYamlAll(str string) ([]any, error) {
	a := []any{}

	d := yaml.NewDecoder(strings.NewReader(str))
	for {
		var v any
		err := d.Decode(&v)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return nil, fmt.Errorf("failed to unmarshal YAML document: %w", err)
		}
		if v != nil {
			a = append(a, v)
		}
	}

	return a, nil
}

// Exec runs external binary and returns its standard output.
// Used as custom template function.
//
// Can be called in 4 ways:
// {{ exec "command arg1 arg2" }} - command string with arguments, split by shell-like rules
// {{ exec "command" (list "arg1" "arg2") }} - command with explicit argument list
// {{ "input" | exec "command arg1 arg2" }} - piped input with command string
// {{ "input" | exec "command" (list "arg1" "arg2") }} - piped input with explicit argument list.
//
//nolint:gocognit,cyclop
func Exec(command string, args ...any) (string, error) {
	var input string
	var strArgs []string

	switch len(args) {
	case 0: // {{ exec "command arg1 arg2" }}
		var err error
		command, strArgs, err = parseCommandArgs(command)
		if err != nil {
			return "", err
		}
	case 1:
		switch v := args[0].(type) {
		case []any: // {{ exec "command" (list "arg1" "arg2") }}
			var err error
			strArgs, err = convertArgsToStrings(v)
			if err != nil {
				return "", err
			}
		case string: // {{ "input" | exec "command arg1 arg2" }}
			input = v
			var err error
			command, strArgs, err = parseCommandArgs(command)
			if err != nil {
				return "", err
			}
		case nil: // {{ exec "command arg1 arg2" }} with nil args
			var err error
			command, strArgs, err = parseCommandArgs(command)
			if err != nil {
				return "", err
			}
		default:
			return "", fmt.Errorf("unexpected type of args[0]: %s", reflect.TypeOf(args[0]))
		}
	case 2: // {{ "input" | exec "command" (list "arg1" "arg2") }}
		argList, ok := args[0].([]any)
		if !ok {
			return "", fmt.Errorf("expected []any for args[0], got %s", reflect.TypeOf(args[0]))
		}
		inputStr, ok := args[1].(string)
		if !ok {
			return "", fmt.Errorf("expected string for args[1] (input), got %s", reflect.TypeOf(args[1]))
		}
		var err error
		strArgs, err = convertArgsToStrings(argList)
		if err != nil {
			return "", err
		}
		input = inputStr
	default:
		return "", fmt.Errorf("exec expects 0-2 arguments after command, got %d", len(args))
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

// convertArgsToStrings converts []any to []string.
func convertArgsToStrings(args []any) ([]string, error) {
	strArgs := make([]string, len(args))
	for i, a := range args {
		switch a := a.(type) {
		case string:
			strArgs[i] = a
		default:
			return nil, fmt.Errorf("unexpected type of arg \"%s\" in args %v at index %d", reflect.TypeOf(a), args, i)
		}
	}

	return strArgs, nil
}

// parseCommandArgs parses a command string into command and arguments.
func parseCommandArgs(command string) (string, []string, error) {
	result, err := shlex.Split(command)
	if err != nil {
		return "", nil, fmt.Errorf("failed to parse command %q: %w", command, err)
	}
	if len(result) == 0 {
		return command, []string{}, nil
	}

	return result[0], result[1:], nil
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
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			// Combine stdout and stderr, only include non-empty parts
			var combinedOutput string
			stdout := strings.TrimSpace(string(bs))
			stderr := strings.TrimSpace(string(exitErr.Stderr))
			switch {
			case stdout != "" && stderr != "":
				combinedOutput = stdout + "\n" + stderr
			case stderr != "":
				combinedOutput = stderr
			case stdout != "":
				combinedOutput = stdout
			}
			if combinedOutput != "" {
				wg.ErrChan() <- fmt.Errorf("command failed: %w: %s", err, combinedOutput)
			} else {
				wg.ErrChan() <- fmt.Errorf("command failed: %w", err)
			}
		} else {
			wg.ErrChan() <- fmt.Errorf("failed to run command: %w", err)
		}

		return
	}

	_, err = output.Write(bs)
	if err != nil {
		wg.ErrChan() <- fmt.Errorf("failed while copying %d bytes from stdout: %w", len(bs), err)
	}
}

// SetValueAtPath sets value in map by dot-separated key path.
// Used as custom template function.
//
//nolint:gocognit,cyclop
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
		case []any:
			idx, err := strconv.Atoi(k)
			if err != nil || idx < 0 || idx >= len(typedCurrent) {
				return nil, fmt.Errorf("failed to walk over path %q: invalid array index %q", path, k)
			}
			current = typedCurrent[idx]
		default:
			return nil, fmt.Errorf(
				"failed to walk over path %q: value for key %q is not a map or array: %v",
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
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(typedCurrent) {
			return nil, fmt.Errorf("failed to set value at path %q: invalid array index %q", path, key)
		}
		typedCurrent[idx] = value
	default:
		return nil, fmt.Errorf(
			"failed to set value at path %q: value for key %q is not a map or array: %v",
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

// GetValueAtPath returns value in map by dot-separated key path.
// First argument is dot-separated key path.
// Second argument is default value if key not found and is optional.
// Third argument is map to search in.
// Used as custom template function.
//
//nolint:gocognit,cyclop
func GetValueAtPath(path string, varArgs ...any) (any, error) {
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
	case []any:
		idx, err := strconv.Atoi(key)
		if err != nil || idx < 0 || idx >= len(typedObj) {
			if defSet {
				return def, nil
			}

			return nil, fmt.Errorf("invalid array index %q for path", key)
		}
		v = typedObj[idx]
	default:
		r, err := tryReflectGet(obj, key, defSet, def)
		if err != nil {
			return nil, err
		}
		v = r
	}

	if defSet {
		return GetValueAtPath(strings.Join(keys[1:], "."), def, v)
	}

	return GetValueAtPath(strings.Join(keys[1:], "."), v)
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

// HasValueAtPath searches for any value by dot-separated key path in map.
// Used as custom template function.
//
//nolint:gocognit,cyclop
func HasValueAtPath(path string, varArgs ...any) (bool, error) {
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
	case []any:
		idx, err := strconv.Atoi(keys[0])
		if err != nil || idx < 0 || idx >= len(typedObj) {
			return defSet, err
		}
		v = typedObj[idx]
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
		return HasValueAtPath(strings.Join(keys[1:], "."), def, v)
	}

	return HasValueAtPath(strings.Join(keys[1:], "."), v)
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
