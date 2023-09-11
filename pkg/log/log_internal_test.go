package log

import (
	"testing"

	"github.com/helmwave/helmwave/pkg/helper"
	"github.com/helmwave/helmwave/pkg/kubedog"
	formatter "github.com/helmwave/logrus-emoji-formatter"
	log "github.com/sirupsen/logrus"
	logTest "github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/suite"
	"github.com/werf/logboek"
	"k8s.io/klog/v2"
)

type LogTestSuite struct {
	suite.Suite

	defaultHooks log.LevelHooks
	logHook      *logTest.Hook
}

func (ts *LogTestSuite) SetupSuite() {
	ts.defaultHooks = log.StandardLogger().Hooks
	ts.logHook = logTest.NewLocal(log.StandardLogger())
}

func (ts *LogTestSuite) TearDownTestSuite() {
	ts.logHook.Reset()
}

func (ts *LogTestSuite) TearDownSuite() {
	log.StandardLogger().ReplaceHooks(ts.defaultHooks)
}

func (ts *LogTestSuite) getLoggerMessages() []string {
	res := make([]string, len(ts.logHook.Entries))
	for i, entry := range ts.logHook.AllEntries() {
		res[i] = entry.Message
	}

	return res
}

func (ts *LogTestSuite) TestKLogHandler() {
	settings := &Settings{
		format: "json",
		level:  "info",
	}
	ts.Require().NoError(settings.Init())

	message := "123"
	klog.Info(message)

	ts.Require().Empty(ts.getLoggerMessages())
}

func (ts *LogTestSuite) TestLogLevel() {
	settings := &Settings{
		format: "text",
		level:  "info",
	}
	ts.Require().NoError(settings.Init())

	log.Debug("test 123")
	ts.Require().Empty(ts.getLoggerMessages(), "message below minimum level should not be logged")
}

func (ts *LogTestSuite) TestDebugLogLevel() {
	settings := &Settings{
		format: "text",
		level:  "debug",
	}
	ts.Require().NoError(settings.Init())

	ts.Require().True(helper.Helm.Debug, "helm debug should be enabled")
}

func (ts *LogTestSuite) TestInvalidLogLevel() {
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
		ts.Require().Error(item.s.Init(), item.msg)
	}
}

func (ts *LogTestSuite) TestFormatter() {
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
		ts.Require().NoError(settings[i].s.Init())
		ts.Require().Equal(settings[i].formatter, log.StandardLogger().Formatter, settings[i].msg)
	}
}

func (ts *LogTestSuite) TestDefaultFormatter() {
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
		ts.Require().NoError(item.s.Init())
		ts.Require().Same(defaultFormatter, log.StandardLogger().Formatter, item.msg)
	}
}

func (ts *LogTestSuite) TestLogboekWidth() {
	width := 1

	kubedog.FixLog(width)
	ts.Require().Equal(width, logboek.DefaultLogger().Streams().Width(), "logboek width should be set")
}

//nolint:paralleltest // helmwave uses single logger for the whole program
func TestLogTestSuite(t *testing.T) {
	// t.Parallel()
	suite.Run(t, new(LogTestSuite))
}
