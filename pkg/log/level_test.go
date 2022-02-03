package log_test

import (
	"testing"

	"github.com/factory-roraimabits/go-deer/pkg/log"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewAtomicLevel(t *testing.T) {
	assert.Equal(t, log.NewAtomicLevel(), zap.NewAtomicLevel())
}

func TestAtomicLevelAt(t *testing.T) {
	testCases := []struct {
		name          string
		level         log.AtomicLevel
		expectedLevel zap.AtomicLevel
	}{
		{
			name:          "Debug",
			level:         log.NewAtomicLevelAt(log.DebugLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.DebugLevel),
		},
		{
			name:          "Info",
			level:         log.NewAtomicLevelAt(log.InfoLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.InfoLevel),
		},
		{
			name:          "Warn",
			level:         log.NewAtomicLevelAt(log.WarnLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.WarnLevel),
		},
		{
			name:          "Error",
			level:         log.NewAtomicLevelAt(log.ErrorLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.ErrorLevel),
		},
		{
			name:          "DPanic",
			level:         log.NewAtomicLevelAt(log.DPanicLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.DPanicLevel),
		},
		{
			name:          "Panic",
			level:         log.NewAtomicLevelAt(log.PanicLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.PanicLevel),
		},
		{
			name:          "Fatal",
			level:         log.NewAtomicLevelAt(log.FatalLevel),
			expectedLevel: zap.NewAtomicLevelAt(zap.FatalLevel),
		},
	}

	for _, tt := range testCases {
		assert.Equal(t, tt.expectedLevel, tt.level, "test of level [%s] failed", tt.name)
	}
}
