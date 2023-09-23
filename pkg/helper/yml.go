package helper

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/helmwave/go-fsimpl"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// SaveInterface encodes input to YAML and saves to file.
func SaveInterface(ctx context.Context, destFS fsimpl.WriteableFS, destFile string, in any) error {
	f, err := CreateFile(destFS, destFile)
	if err != nil {
		return err
	}

	defer func(f fs.File) {
		err := f.Close()
		if err != nil {
			log.WithError(err).Error("failed to close file")
		}
	}(f)

	data := Byte(ctx, in)

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to %s: %w", destFile, err)
	}

	return nil
}

// Byte marshals input to YAML and returns YAML byte slice.
func Byte(ctx context.Context, in any) []byte {
	data, err := yaml.Marshal(in)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
