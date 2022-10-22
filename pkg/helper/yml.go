package helper

import (
	"context"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
	log "github.com/sirupsen/logrus"
)

// SaveInterface encodes input to YAML and saves to file.
func SaveInterface(ctx context.Context, file string, in interface{}) error {
	f, err := CreateFile(file)
	if err != nil {
		return err
	}

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("failed to close file %s: %v", f.Name(), err)
		}
	}(f)

	data := Byte(ctx, in)

	_, err = f.Write(data)
	if err != nil {
		return fmt.Errorf("failed to marshal to %s: %w", file, err)
	}

	if err = f.Sync(); err != nil {
		return fmt.Errorf("failed to sync file %s: %w", file, err)
	}

	return nil
}

// Byte marshals input to YAML and returns YAML byte slice.
func Byte(ctx context.Context, in interface{}) []byte {
	data, err := yaml.MarshalContext(ctx, in)
	if err != nil {
		log.Fatal(err)
	}

	return data
}
