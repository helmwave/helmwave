package helper

import (
	"fmt"
	"os"

	"github.com/helmwave/go-fsimpl"
	dir "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

// MoveFile moves files or directories. It also handles move between different mounts (copy + rm).
func MoveFile(srcFS, dstFS fsimpl.WriteableFS, src, dst string) error {
	// TODO: use srcFS and dstFS
	// It doesnt work if workdir has been mounted.
	err := os.Rename(src, dst)
	if err != nil {
		err = dir.Copy(src, dst)
		if err != nil {
			return fmt.Errorf("failed to move file between filesystems: %w", err)
		}
		defer func(src string) {
			err := os.RemoveAll(src)
			if err != nil {
				log.WithError(err).Error("failed to remove source temporary directory")
			}
		}(src)
	}

	return nil
}
