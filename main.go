package main

import (
	"bytes"
	"fmt"
	"github.com/helmwave/helmwave/pkg/plan"
	log "github.com/sirupsen/logrus"
	_ "go.beyondstorage.io/services/fs/v4"
	"go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/services"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

func main() {
	b := plan.N()

	path := "tests"
	file := "02_helmwave.yml"

	src := "fs://" + path

	buf := new(bytes.Buffer)
	cudir, _ := os.Getwd()

	store, err := services.NewStoragerFromString(src, pairs.WithWorkDir(cudir))
	if err != nil {
		log.Fatal("init err: ", err)
	}
	_, err = store.Read(filepath.Join(path, file), buf)
	if err != nil {
		log.Fatal("read err: ", err)
	}

	err = yaml.Unmarshal(buf.Bytes(), b)

	fmt.Println(b)
}
