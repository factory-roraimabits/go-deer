package log_test

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"

	"fmt"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/factory-roraimabits/go-deer/pkg/log"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// logRegex returns the level as the first group, discards the timestamp, logger as the
	// second group, caller is discarded and everything after that as the fourth group.
	//
	// Examples of matched lines:
	//   [ts:2019-04-01T15:39:09.142773Z][level:debug][caller:log/logger_test.go:21][msg:before contextualising]
	//   [ts:2019-04-01T17:19:16.290081Z][level:warn][logger:first_level][caller:log/logger_test.go:97][msg:message]
	logRegex = regexp.MustCompile(`\[ts:(?:[0-9-T:.]+Z)]\[level:([a-z]+)](\[logger:(?:.*?)])?\[caller:(.*?)](.*)`)

	// stacktraceRegex finds the stacktrace segment within a log line.
	stacktraceRegex = regexp.MustCompile(`(\[stacktrace:(?:.*?)])`)
)

type LogLine struct {
	Level      string
	LoggerName string
	Message    string
}

func TestKeyValueLogger(t *testing.T) {
	parseLogLine := func(t *testing.T, line string) LogLine {
		matches := logRegex.FindAllStringSubmatch(line, -1)

		if len(matches[0]) != 5 {
			t.Fatalf("expected regex to have 5 matches, %d found", len(matches[0]))
		}

		lvl, name, msg := matches[0][1], matches[0][2], matches[0][4]

		return LogLine{
			Level:      lvl,
			LoggerName: name,
			Message:    msg,
		}
	}

	assertLine := func(t *testing.T, line, level, content, name string) {
		l := parseLogLine(t, line)

		if l.Level != level {
			t.Fatalf("expected log level to be %s, got: %s", level, l.Level)
		}

		if l.Message != content {
			t.Fatalf("expected content to be %s, got: %s", content, l.Message)
		}

		if l.LoggerName != name {
			t.Fatalf("expected logger name to be %s, got: %s", name, l.LoggerName)
		}
	}

	assertAndRemoveStacktrace := func(t *testing.T, lines []string) string {
		l := strings.Join(lines, "")
		if !stacktraceRegex.MatchString(l) {
			t.Fatalf("expected line to have stacktrace, none found")
		}
		return stacktraceRegex.ReplaceAllString(l, "")
	}

	tt := []struct {
		Name       string
		Level      zapcore.Level
		SetupFunc  func(t *testing.T, l log.Logger)
		AssertFunc func(t *testing.T, lines []string)
	}{
		{
			Name:  "Log Using Raw Logger",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				l.Debug("my Debug message")
				l.Info("my Info message")
				l.Warn("my Warn message")
				l.Error("my Error message")
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message]", "")
				assertLine(t, lines[1], "info", "[msg:my Info message]", "")
				assertLine(t, lines[2], "warn", "[msg:my Warn message]", "")
				assertLine(t, assertAndRemoveStacktrace(t, lines[3:]), "error", `[msg:my Error message]`, "")
			},
		},
		{
			Name:  "Log Using Context",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)

				log.Debug(ctx, "my Debug message")
				log.Info(ctx, "my Info message")
				log.Warn(ctx, "my Warn message")
				log.Error(ctx, "my Error message")
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message]", "")
				assertLine(t, lines[1], "info", "[msg:my Info message]", "")
				assertLine(t, lines[2], "warn", "[msg:my Warn message]", "")
				assertLine(t, assertAndRemoveStacktrace(t, lines[3:]), "error", `[msg:my Error message]`, "")
			},
		},
		{
			Name:  "Named Logger",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)

				ctx = log.Named(ctx, "first_level")
				log.Debug(ctx, "my Debug message")

				ctx = log.Named(ctx, "second_level")
				log.Info(ctx, "my Info message")

				ctx = log.Named(ctx, "third_level")
				log.Warn(ctx, "my Warn message")
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message]", "[logger:first_level]")
				assertLine(t, lines[1], "info", "[msg:my Info message]", "[logger:first_level.second_level]")
				assertLine(t, lines[2], "warn", "[msg:my Warn message]", "[logger:first_level.second_level.third_level]")
			},
		},
		{
			Name:  "Check Works (Should log)",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)
				if ce := log.Check(ctx, zap.DebugLevel, "my Debug message"); ce != nil {
					ce.Write(zap.String("foo", "bar"))
				}
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message][foo:bar]", "")
			},
		},
		{
			Name:  "Check Works (Should not log)",
			Level: zap.InfoLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)
				if ce := log.Check(ctx, zap.DebugLevel, "my Debug message"); ce != nil {
					ce.Write(zap.String("foo", "bar"))
				}
			},
		},
		{
			Name:  "Log Message With Fields",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)

				log.Debug(ctx, "my Debug message",
					zap.String("string_key", "value"),
					zap.Time("time_key", time.Unix(0, 0)),
					zap.Int64("int64_key", 1234),
					zap.Float64("float64_key", 1234.5678),
					zap.Error(fmt.Errorf("my error")),
				)
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message][string_key:value][time_key:1970-01-01T00:00:00.000000Z][int64_key:1234][float64_key:1234.5678][error:my error]", "") // nolint
			},
		},
		{
			Name:  "Log Message With Context Fields",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)

				ctx = log.With(ctx,
					zap.String("string_key", "value"),
					zap.Time("time_key", time.Unix(0, 0)),
					zap.Int64("int64_key", 1234),
					zap.Float64("float64_key", 1234.5678),
					zap.Error(fmt.Errorf("my error")),
					zap.Duration("duration_key", 374*time.Millisecond),
				)

				log.Debug(ctx, "my Debug message", zap.String("extra", "debug_extra"))
				log.Info(ctx, "my Info message", zap.String("extra", "info_extra"))
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message][string_key:value][time_key:1970-01-01T00:00:00.000000Z][int64_key:1234][float64_key:1234.5678][error:my error][duration_key:0.374][extra:debug_extra]", "") // nolint
				assertLine(t, lines[1], "info", "[msg:my Info message][string_key:value][time_key:1970-01-01T00:00:00.000000Z][int64_key:1234][float64_key:1234.5678][error:my error][duration_key:0.374][extra:info_extra]", "")    // nolint
			},
		},
		{
			Name:  "Test Panic Levels",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				defer func() {
					if r := recover(); r == nil {
						t.Fatal("expected panic to happen")
					}
				}()

				ctx := log.Context(context.Background(), l)
				log.Panic(ctx, "my Panic message")
			},
			AssertFunc: func(t *testing.T, lines []string) {
				line := assertAndRemoveStacktrace(t, lines)
				assertLine(t, line, "panic", "[msg:my Panic message]", "")
			},
		},
		{
			Name:  "Test DPanic Levels",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)
				log.DPanic(ctx, "my DPanic message")
			},
			AssertFunc: func(t *testing.T, lines []string) {
				line := assertAndRemoveStacktrace(t, lines)
				assertLine(t, line, "dpanic", "[msg:my DPanic message]", "")
			},
		},
		{
			Name:  "Test Sugar Logger",
			Level: zap.DebugLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx := log.Context(context.Background(), l)

				logger := log.Sugar(ctx)
				logger.Debugw("my Debug message", "string_key", "value", "int64_key", 123456)
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "debug", "[msg:my Debug message][string_key:value][int64_key:123456]", "")
			},
		},
		{
			Name:  "Test Change Levels",
			Level: zap.ErrorLevel,
			SetupFunc: func(t *testing.T, l log.Logger) {
				ctx1 := log.Context(context.Background(), l)
				log.Debug(ctx1, "should not appear", zap.String("log_level", "error"))
				log.Info(ctx1, "should not appear", zap.String("log_level", "error"))

				ctx2 := log.WithLevel(ctx1, zap.InfoLevel)
				// Previous contexts should remain at their own level.
				log.Debug(ctx1, "should not appear", zap.String("log_level", "error"))
				log.Info(ctx1, "should not appear", zap.String("log_level", "error"))

				// New context should accept new level.
				log.Debug(ctx2, "should not appear", zap.String("log_level", "info"))
				log.Info(ctx2, "should appear", zap.String("log_level", "info"))

				ctx3 := log.WithLevel(ctx2, zap.DebugLevel)
				// Previous contexts should remain at their own level.
				log.Debug(ctx1, "should not appear", zap.String("log_level", "error"))
				log.Info(ctx1, "should not appear", zap.String("log_level", "error"))
				log.Debug(ctx2, "should not appear", zap.String("log_level", "info"))

				// New context should accept new level.
				log.Debug(ctx3, "should appear", zap.String("log_level", "debug"))
				log.Info(ctx3, "should appear", zap.String("log_level", "debug"))
			},
			AssertFunc: func(t *testing.T, lines []string) {
				assertLine(t, lines[0], "info", "[msg:should appear][log_level:info]", "")
				assertLine(t, lines[1], "debug", "[msg:should appear][log_level:debug]", "")
				assertLine(t, lines[2], "info", "[msg:should appear][log_level:debug]", "")
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			var out bytes.Buffer

			lvl := zap.NewAtomicLevelAt(tc.Level)
			l := log.NewProductionLogger(&lvl, log.WithWriter(zapcore.AddSync(&out)))
			tc.SetupFunc(t, l)

			var lines []string

			s := bufio.NewScanner(&out)
			for s.Scan() {
				lines = append(lines, s.Text())
			}

			if err := s.Err(); err != nil {
				t.Fatalf("error reading stdErr output buffer: %v", err)
			}

			// If there is no AssertFunc then we expect no log lines.
			if tc.AssertFunc == nil {
				require.Zero(t, lines)
				return
			}

			require.NotZero(t, lines)
			tc.AssertFunc(t, lines)
		})
	}
}

func TestJSONLogger(t *testing.T) {
	var out bytes.Buffer

	lvl := zap.NewAtomicLevelAt(log.DebugLevel)
	l := log.NewProductionLogger(&lvl,
		log.WithWriter(zapcore.AddSync(&out)),
		log.WithJSONEncoding(),
	)

	l.Debug("my long message", log.String("custom_tag", "my value"), log.Bool("boolean", true))

	msg := struct {
		Boolean   bool      `json:"boolean"`
		Caller    string    `json:"caller"`
		CustomTag string    `json:"custom_tag"`
		Level     string    `json:"level"`
		Message   string    `json:"msg"`
		Timestamp time.Time `json:"ts"`
	}{}

	err := json.NewDecoder(&out).Decode(&msg)
	require.NoError(t, err)

	require.Equal(t, "my long message", msg.Message)
	require.Equal(t, "my value", msg.CustomTag)
	require.Equal(t, true, msg.Boolean)
	require.Equal(t, "debug", msg.Level)
}

func TestConsoleLogger(t *testing.T) {
	var out bytes.Buffer

	lvl := zap.NewAtomicLevelAt(log.DebugLevel)
	l := log.NewProductionLogger(&lvl,
		log.WithWriter(zapcore.AddSync(&out)),
		log.WithConsoleEncoding(),
	)

	l.Debug("my long message", log.String("custom_tag", "my value"), log.Bool("boolean", true))
	require.True(t, strings.Contains(out.String(), `my long message	{"custom_tag": "my value", "boolean": true}`))
}

func TestFromContext(t *testing.T) {
	logger := log.FromContext(context.Background())
	require.Nil(t, logger)

	lvl := zap.NewAtomicLevelAt(zapcore.DebugLevel)
	l := log.NewProductionLogger(&lvl)

	ctx := log.Context(context.Background(), l)
	logger = log.FromContext(ctx)

	require.Same(t, l, logger)
}

func TestLoggerOptions(t *testing.T) {
	type Message struct {
		Caller     string `json:"caller"`
		Level      string `json:"level"`
		LogLevel   string `json:"log_level"`
		Stacktrace string `json:"stacktrace"`
	}

	setup := func(t *testing.T, opt log.Option) Message {
		var out bytes.Buffer

		lvl := zap.NewAtomicLevelAt(log.DebugLevel)
		l := log.NewProductionLogger(&lvl,
			log.WithWriter(zapcore.AddSync(&out)),
			log.WithJSONEncoding(),
			opt,
		)

		l.Error("my long message")

		var msg Message
		err := json.NewDecoder(&out).Decode(&msg)
		require.NoError(t, err)
		return msg
	}

	t.Run("Without Caller", func(t *testing.T) {
		msg := setup(t, log.WithCaller(false))
		require.Empty(t, msg.Caller)
	})

	t.Run("Without Stacktrace", func(t *testing.T) {
		msg := setup(t, log.WithStacktraceOnError(false))
		require.Empty(t, msg.Stacktrace)
	})

	t.Run("With Custom Level Key", func(t *testing.T) {
		msg := setup(t, log.WithLevelKey("log_level"))
		require.Empty(t, msg.Level)
		require.Equal(t, "error", msg.LogLevel)
	})
}
