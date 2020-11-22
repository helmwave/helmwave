package yml

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func Read(file string, yml *Body) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}

	err = yaml.Unmarshal(src, yml)
	if err != nil {
		log.Fatal(err)
	}
}

func Save(file string, in interface{}) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	data := Byte(in)

	f.Write(data)
	return f.Close()

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
