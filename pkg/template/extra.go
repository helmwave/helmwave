package template

import (
	"fmt"
	"golang.org/x/sync/errgroup"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

type Values = map[string]interface{}

func ToYaml(v interface{}) (string, error) {
	data, err := yaml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func FromYaml(str string) (Values, error) {
	m := Values{}

	if err := yaml.Unmarshal([]byte(str), &m); err != nil {
		return nil, fmt.Errorf("%s, offending yaml: %s", err, str)
	}
	return m, nil
}

func Exec(command string, args []interface{}, inputs ...string) (string, error) {
	var input string
	if len(inputs) > 0 {
		input = inputs[0]
	}

	strArgs := make([]string, len(args))
	for i, a := range args {
		switch a.(type) {
		case string:
			strArgs[i] = a.(string)
		default:
			return "", fmt.Errorf("unexpected type of arg \"%s\" in args %v at index %d", reflect.TypeOf(a), args, i)
		}
	}

	cmd := exec.Command(command, strArgs...)
	//cmd.Dir = c.basePath

	g := errgroup.Group{}

	if len(input) > 0 {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return "", err
		}

		g.Go(func() error {
			defer stdin.Close()

			size := len(input)

			i := 0

			for {
				n, err := io.WriteString(stdin, input[i:])
				if err != nil {
					return fmt.Errorf("failed while writing %d bytes to stdin of \"%s\": %v", len(input), command, err)
				}

				i += n

				if i == size {
					return nil
				}
			}
		})
	}

	var bytes []byte

	g.Go(func() error {
		bs, err := cmd.Output()
		if err != nil {
			return err
		}

		bytes = bs
		return nil
	})

	if err := g.Wait(); err != nil {
		return "", err
	}

	return string(bytes), nil
}

func SetValueAtPath(path string, value interface{}, values Values) (Values, error) {
	var current interface{}
	current = values
	components := strings.Split(path, ".")
	pathToMap := components[:len(components)-1]
	key := components[len(components)-1]
	for _, k := range pathToMap {
		var elem interface{}

		switch typedCurrent := current.(type) {
		case map[string]interface{}:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" does not exist", path, k)
			}
			elem = v
		case map[interface{}]interface{}:
			v, exists := typedCurrent[k]
			if !exists {
				return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" does not exist", path, k)
			}
			elem = v
		default:
			return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, k)
		}

		switch typedElem := elem.(type) {
		case map[string]interface{}, map[interface{}]interface{}:
			current = typedElem
		default:
			return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, k)
		}
	}

	switch typedCurrent := current.(type) {
	case map[string]interface{}:
		typedCurrent[key] = value
	case map[interface{}]interface{}:
		typedCurrent[key] = value
	default:
		return nil, fmt.Errorf("failed to set value at path \"%s\": value for key \"%s\" was not a map", path, key)
	}
	return values, nil
}

func RequiredEnv(name string) (string, error) {
	if val, exists := os.LookupEnv(name); exists && len(val) > 0 {
		return val, nil
	}

	return "", fmt.Errorf("required env var `%s` is not set", name)
}

func Required(warn string, val interface{}) (interface{}, error) {
	if val == nil {
		return nil, fmt.Errorf(warn)
	} else if _, ok := val.(string); ok {
		if val == "" {
			return nil, fmt.Errorf(warn)
		}
	}

	return val, nil
}

func ReadFile(file string) (string, error) {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

type noValueError struct {
	msg string
}

func (e *noValueError) Error() string {
	return e.msg
}

func Get(path string, varArgs ...interface{}) (interface{}, error) {
	var defSet bool
	var def interface{}
	var obj interface{}
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
		return nil, fmt.Errorf("unexpected number of args pased to the template function get(path, [def, ]obj): expected 1 or 2, got %d, args was %v", len(varArgs), varArgs)
	}

	if path == "" {
		return obj, nil
	}
	keys := strings.Split(path, ".")
	var v interface{}
	var ok bool
	switch typedObj := obj.(type) {
	case map[string]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no value exist for key \"%s\" in %v", keys[0], typedObj)}
		}
	case map[interface{}]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no value exist for key \"%s\" in %v", keys[0], typedObj)}
		}
	default:
		maybeStruct := reflect.ValueOf(typedObj)
		if maybeStruct.Kind() != reflect.Struct {
			return nil, &noValueError{fmt.Sprintf("unexpected type(%v) of value for key \"%s\": it must be either map[string]interface{} or any struct", reflect.TypeOf(obj), keys[0])}
		} else if maybeStruct.NumField() < 1 {
			return nil, &noValueError{fmt.Sprintf("no accessible struct fields for key \"%s\"", keys[0])}
		}
		f := maybeStruct.FieldByName(keys[0])
		if !f.IsValid() {
			if defSet {
				return def, nil
			}
			return nil, &noValueError{fmt.Sprintf("no field named \"%s\" exist in %v", keys[0], typedObj)}
		}
		v = f.Interface()
	}

	if defSet {
		return Get(strings.Join(keys[1:], "."), def, v)
	}
	return Get(strings.Join(keys[1:], "."), v)
}

func HasKey(path string, varArgs ...interface{}) (bool, error) {
	var defSet bool
	var def interface{}
	var obj interface{}
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
		return false, fmt.Errorf("unexpected number of args pased to the template function get(path, [def, ]obj): expected 1 or 2, got %d, args was %v", len(varArgs), varArgs)
	}

	if path == "" {
		return true, nil
	}
	keys := strings.Split(path, ".")
	var v interface{}
	var ok bool
	switch typedObj := obj.(type) {
	case map[string]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			return defSet, nil
		}
	case map[interface{}]interface{}:
		v, ok = typedObj[keys[0]]
		if !ok {
			return defSet, nil
		}
	default:
		maybeStruct := reflect.ValueOf(typedObj)
		if maybeStruct.Kind() != reflect.Struct {
			return false, &noValueError{fmt.Sprintf("unexpected type(%v) of value for key \"%s\": it must be either map[string]interface{} or any struct", reflect.TypeOf(obj), keys[0])}
		} else if maybeStruct.NumField() < 1 {
			return false, &noValueError{fmt.Sprintf("no accessible struct fields for key \"%s\"", keys[0])}
		}
		f := maybeStruct.FieldByName(keys[0])
		if !f.IsValid() {
			return defSet, nil
		}
		v = f.Interface()
	}

	if defSet {
		return HasKey(strings.Join(keys[1:], "."), def, v)
	}
	return HasKey(strings.Join(keys[1:], "."), v)
}
