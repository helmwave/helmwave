package yml

import (
	"github.com/helmwave/helmwave/pkg/helper"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Read(file string, yml *Config) error {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(src, yml)
	if err != nil {
		return err
	}

	return nil
}

func Save(file string, in interface{}) error {
	f, err := helper.CreateFile(file)
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

func Print(in interface{}) {
	println(String(in))
}
