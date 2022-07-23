package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mgutz/ansi"
)

const (
	startSectionFormat          = ansi.Reset + "section_start:%v:%s\r" + ansi.Reset + "%s\n"
	startCollapsedSectionFormat = ansi.Reset + "section_start:%v:%s[collapsed=true]\r" + ansi.Reset + "%s\n"
	stopSectionFormat           = ansi.Reset + "section_end:%v:%s\r" + ansi.Reset + "\n"
)

var ErrNotGitlabCI = fmt.Errorf("Current environment is not Gitlab CI (GITLAB_CI!=true)")

type GitlabLogSection struct {
	Name      string
	Header    string
	collapsed bool
	contents  *bytes.Buffer

	startTimestamp time.Time
}

func NewSection(name string) (GitlabLogSection, error) {
	if v, defined := os.LookupEnv("GITLAB_CI"); !defined || (v != "true") {
		return GitlabLogSection{}, ErrNotGitlabCI
	}

	i := GitlabLogSection{
		Name:           name,
		Header:         name,
		startTimestamp: time.Now(),
		contents:       &bytes.Buffer{},
	}

	return i, i.startSection()
}

func (section GitlabLogSection) Write(p []byte) (int, error) {
	return section.contents.Write(p)
}

func (section GitlabLogSection) WriteString(s string) (int, error) {
	return section.contents.WriteString(s)
}

func (section GitlabLogSection) startSection() error {
	format := startSectionFormat
	if section.collapsed {
		format = startCollapsedSectionFormat
	}

	return section.writeControlSectionCommand(format)
}

func (section GitlabLogSection) stopSection() error {
	return section.writeControlSectionCommand(stopSectionFormat)
}

func (section GitlabLogSection) writeControlSectionCommand(format string) error {
	_, err := fmt.Fprintf(section.contents, format, time.Now().Unix(), section.Name, section.Name)

	return err
}

func (section GitlabLogSection) Fire(writer io.Writer) error {
	section.stopSection()
	_, err := section.contents.WriteTo(writer)

	return err
}
