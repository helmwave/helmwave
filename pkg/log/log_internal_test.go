package log

import (
	"bytes"
	"os"
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	formatter "github.com/helmwave/logrus-emoji-formatter"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/werf/logboek"
	"k8s.io/klog/v2"
)

var buf bytes.Buffer

type LogTestSuite struct {
	suite.Suite
}

func (s *LogTestSuite) SetupTest() {
	log.StandardLogger().SetOutput(&buf)
}

func (s *LogTestSuite) TearDownTest() {
	log.StandardLogger().SetOutput(os.Stderr)
}

func (s *LogTestSuite) TestKLogHandler() {
	settings := &Settings{
		format: "json",
		level:  "info",
	}
	s.Require().NoError(settings.Init())

	message := "123"
	klog.Info(message)

	s.Require().Zero(buf.Len())
	buf.Reset()
}

func (s *LogTestSuite) TestLogLevel() {
	settings := &Settings{
		format: "text",
		level:  "info",
	}
	s.Require().NoError(settings.Init())

	log.Debug("test 123")
	defer buf.Reset()
	s.Require().Empty(buf.String(), "message below minimum level should not be logged")
}

func (s *LogTestSuite) TestDebugLogLevel() {
	settings := &Settings{
		format: "text",
		level:  "debug",
	}
	s.Require().NoError(settings.Init())

	s.Require().True(helper.Helm.Debug, "helm debug should be enabled")
}

func (s *LogTestSuite) TestInvalidLogLevel() {
	settings := []struct {
		s   *Settings
		msg string
	}{
		{
			s: &Settings{
				format: "text",
			},
			msg: "should error with no level",
		},
		{
			s: &Settings{
				format: "text",
				level:  "blabla123",
			},
			msg: "should error with invalid level",
		},
	}

	for _, item := range settings {
		s.Require().Error(item.s.Init(), item.msg)
	}
}

func (s *LogTestSuite) TestFormatter() {
	settings := []struct {
		s         *Settings
		formatter log.Formatter
		msg       string
	}{
		{
			s: &Settings{
				format: "json",
				level:  "info",
			},
			formatter: &log.JSONFormatter{
				PrettyPrint: true,
			},
			msg: "should use json formatter",
		},
		{
			s: &Settings{
				format: "pad",
				level:  "info",
			},
			formatter: &log.TextFormatter{
				PadLevelText:     true,
				DisableTimestamp: true,
			},
			msg: "should use pad formatter",
		},
		{
			s: &Settings{
				format: "emoji",
				level:  "info",
			},
			formatter: &formatter.Config{
				LogFormat: "[%lvl%]: %msg%",
			},
			msg: "should use emoji formatter",
		},
		{
			s: &Settings{
				format: "text",
				level:  "info",
			},
			formatter: &log.TextFormatter{
				DisableTimestamp: true,
			},
			msg: "should use text formatter",
		},
	}

	for i := range settings {
		s.Require().NoError(settings[i].s.Init())
		s.Require().Equal(settings[i].formatter, log.StandardLogger().Formatter, settings[i].msg)
	}
}

func (s *LogTestSuite) TestDefaultFormatter() {
	defaultFormatter := &log.TextFormatter{}
	log.SetFormatter(defaultFormatter)

	settings := []struct {
		s   *Settings
		msg string
	}{
		{
			s: &Settings{
				level: "info",
			},
			msg: "should use default formatter",
		},
		{
			s: &Settings{
				format: "blabla123",
				level:  "info",
			},
			msg: "should use default formatter",
		},
	}

	for _, item := range settings {
		s.Require().NoError(item.s.Init())
		s.Require().Same(defaultFormatter, log.StandardLogger().Formatter, item.msg)
	}
}

func (s *LogTestSuite) TestLogboekWidth() {
	settings := &Settings{
		level:  "info",
		format: "text",
		width:  1,
	}

	s.Require().NoError(settings.Init())
	s.Require().NoError(kubedog.FixKubedogLog(settings.width))
	s.Require().Equal(settings.width, logboek.DefaultLogger().Streams().Width(), "logboek width should be set")
}

func TestLogTestSuite(t *testing.T) { //nolintlint:paralleltest // helmwave uses single logger for the whole program
	// t.Parallel()
	suite.Run(t, new(LogTestSuite))
}
