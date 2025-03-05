package audit

import (
	"context"
	"encoding/json"
	"log/slog"
)

// Auditor log handler
func (a *Auditor) Handle(ctx context.Context, r slog.Record) error {
	data := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.Any()
		return true
	})

	event := Event{
		Time:    r.Time.Format(cfg.TimeFormat),
		Level:   r.Level.String(),
		Message: r.Message,
		Data:    data,
	}

	jsonEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}

	a.Logger.Println(string(jsonEvent))

	return nil
}
