package tracker

import (
	"flag"
	"fmt"
	"io"
	"os"

	klog "k8s.io/klog/v2"
)

// SilenceKlog discards all klog logs except FATAL: the client-go machinery behind the watcher is
// chatty on levels that mean nothing to a helmwave user.
func SilenceKlog() error {
	fs := flag.NewFlagSet("klog", flag.PanicOnError)
	klog.InitFlags(fs)

	if err := fs.Set("logtostderr", "false"); err != nil {
		return fmt.Errorf("failed to disable 'logtostderr': %w", err)
	}
	if err := fs.Set("alsologtostderr", "false"); err != nil {
		return fmt.Errorf("failed to disable 'alsologtostderr': %w", err)
	}
	if err := fs.Set("stderrthreshold", "5"); err != nil {
		return fmt.Errorf("failed to disable 'stderrthreshold': %w", err)
	}

	klog.SetOutputBySeverity("INFO", io.Discard)
	klog.SetOutputBySeverity("WARNING", io.Discard)
	klog.SetOutputBySeverity("ERROR", io.Discard)
	klog.SetOutputBySeverity("FATAL", os.Stderr)

	return nil
}
