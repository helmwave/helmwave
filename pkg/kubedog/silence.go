package kubedog

import (
	"context"
	"flag"
	"fmt"
	"io"

	"github.com/werf/logboek"
	"k8s.io/klog"
	klogV2 "k8s.io/klog/v2"
)

// SilenceKlogV2 discards all klog/v2 logs except FATAL.
func SilenceKlogV2(ctx context.Context) error {
	fs := flag.NewFlagSet("klog", flag.PanicOnError)
	klogV2.InitFlags(fs)

	if err := silenceKlogFlagSet(fs); err != nil {
		return err
	}

	// Suppress info and warnings from client-go reflector
	klogV2.SetOutputBySeverity("INFO", io.Discard)
	klogV2.SetOutputBySeverity("WARNING", io.Discard)
	klogV2.SetOutputBySeverity("ERROR", io.Discard)
	klogV2.SetOutputBySeverity("FATAL", logboek.Context(ctx).ErrStream())

	return nil
}

// SilenceKlog discards all klog logs except FATAL.
func SilenceKlog(ctx context.Context) error {
	fs := flag.NewFlagSet("klog", flag.PanicOnError)
	klog.InitFlags(fs)

	if err := silenceKlogFlagSet(fs); err != nil {
		return err
	}

	// Suppress info and warnings from client-go reflector
	klog.SetOutputBySeverity("INFO", io.Discard)
	klog.SetOutputBySeverity("WARNING", io.Discard)
	klog.SetOutputBySeverity("ERROR", io.Discard)
	klog.SetOutputBySeverity("FATAL", logboek.Context(ctx).ErrStream())

	return nil
}

func silenceKlogFlagSet(fs *flag.FlagSet) error {
	if err := fs.Set("logtostderr", "false"); err != nil {
		return fmt.Errorf("failed to disable 'logtostderr': %w", err)
	}
	if err := fs.Set("alsologtostderr", "false"); err != nil {
		return fmt.Errorf("failed to disable 'alsologtostderr': %w", err)
	}
	if err := fs.Set("stderrthreshold", "5"); err != nil {
		return fmt.Errorf("failed to disable 'stderrthreshold': %w", err)
	}

	return nil
}

// FixKubedogLog will disable kubernetes logger and fix width for logboek.
// Todo: add ctx as an argument.
func FixKubedogLog(width int) error {
	if err := SilenceKlog(context.Background()); err != nil {
		return err
	}

	if err := SilenceKlogV2(context.Background()); err != nil {
		return err
	}

	if width > 0 {
		logboek.DefaultLogger().Streams().SetWidth(width)
	}

	return nil
}
