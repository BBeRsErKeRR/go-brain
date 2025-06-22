package logger

import (
	"bytes"
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// https://github.com/uber-go/zap/issues/829
func init() {
	zap.RegisterEncoder("prettier", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		jsonConfig := cfg
		jsonConfig.TimeKey = ""
		jsonConfig.LevelKey = ""
		jsonConfig.NameKey = ""
		jsonConfig.CallerKey = ""
		jsonConfig.MessageKey = ""
		jsonConfig.StacktraceKey = ""

		consoleConfig := cfg
		consoleConfig.StacktraceKey = ""

		pce := prettyConsoleEncoder{
			EncoderConfig:  &cfg,
			Encoder:        zapcore.NewConsoleEncoder(jsonConfig),
			consoleEncoder: zapcore.NewConsoleEncoder(consoleConfig),
		}

		return &pce, nil
	})
}

type prettyConsoleEncoder struct {
	*zapcore.EncoderConfig
	zapcore.Encoder
	consoleEncoder zapcore.Encoder
}

func (pce *prettyConsoleEncoder) Clone() zapcore.Encoder {
	return &prettyConsoleEncoder{
		EncoderConfig:  pce.EncoderConfig,
		Encoder:        pce.Encoder.Clone(),
		consoleEncoder: pce.consoleEncoder.Clone(),
	}
}

func (pce *prettyConsoleEncoder) EncodeEntry(ent zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	line, err := pce.consoleEncoder.EncodeEntry(ent, nil)
	if err != nil {
		return line, err
	}
	if len(fields) > 0 {
		jsonRawBuf, err := pce.Encoder.EncodeEntry(ent, fields)
		if err != nil {
			return line, err
		}
		defer jsonRawBuf.Free()

		jsonBuf := buffer.NewPool().Get()
		defer jsonBuf.Free()

		var out bytes.Buffer
		err = json.Indent(&out, jsonRawBuf.Bytes(), "", "  ")
		if err != nil {
			return line, err
		}
		_, err = line.Write(out.Bytes())
		if err != nil {
			return line, err
		}
	}

	// If there's no stacktrace key, honor that; this allows users to force
	// single-line output.
	if ent.Stack != "" && pce.StacktraceKey != "" {
		line.AppendByte('\n')
		line.AppendString(ent.Stack)
	}

	// if pce.LineEnding != "" {
	// 	line.AppendString(pce.LineEnding)
	// } else {
	// 	line.AppendString(zapcore.DefaultLineEnding)
	// }
	return line, err
}
