package helper

import (
	"fmt"
	"os"

	dir "github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

// MoveFile moves files or directories. It also handles move between different mounts (copy + rm).
func MoveFile(src, dst string) error {
	log.WithFields(log.Fields{
		"src": src,
		"dst": dst,
	}).Trace("Moving")

	// It doesn't work if workdir has been mounted.
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
