package helper

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func SaveInterface(file string, in interface{}) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}
	defer f.Close()

	data := Byte(in)

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return f.Sync()
}

func Byte(in interface{}) []byte {
	data, err := yaml.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}

	return data
}

func String(in interface{}) string {
	return string(Byte(in))
}
