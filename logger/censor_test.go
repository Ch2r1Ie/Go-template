package logger

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"os"
	"testing"
)

func TestCensorReplacer(t *testing.T) {
	t.Run("censor replacer", func(t *testing.T) {
		os.Unsetenv("ENV")
		buf := bytes.NewBuffer([]byte{})

		defaultLogOutput = buf
		defer func() { defaultLogOutput = os.Stdout }()

		censors = map[string]string{
			"cid": "xxxxxxxxxxxxx",
		}

		l := New(CensorReplacer)

		l.Info("message", slog.String("cid", "123456789012"))

		var m map[string]string
		json.NewDecoder(buf).Decode(&m)
		v, ok := m["cid"]
		if !ok {
			t.Errorf("not found cid key\n")
		}

		if v != "xxxxxxxxxxxxx" {
			t.Errorf("replacer replace cid to %q: actual %q\n", "xxxxxxxxxxxxx", v)
		}
	})
}
