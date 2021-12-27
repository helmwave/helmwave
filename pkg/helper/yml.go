package helper

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SaveInterface encodes input to YAML and saves to file.
func SaveInterface(file string, in interface{}) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck // TODO: need to check error

	data := Byte(in)

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}

// Byte marshals input to YAML and returns YAML byte slice.
func Byte(in interface{}) []byte {
	data, err := yaml.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

// String marshals input to YAML and returns YAML string.
func String(in interface{}) string {
	return string(Byte(in))
}
