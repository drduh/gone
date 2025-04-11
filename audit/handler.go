package audit

import (
	"context"
	"encoding/json"
	"log/slog"
)

// Handle formats and outputs audit events in JSON format.
func (a *Auditor) Handle(ctx context.Context, r slog.Record) error {
	data := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(attr slog.Attr) bool {
		data[attr.Key] = attr.Value.Any()
		return true
	})

	event, err := json.Marshal(&Event{
		Time:    r.Time.Format(a.TimeFormat),
		Level:   r.Level.String(),
		Message: r.Message,
		Data:    data,
	})
	if err != nil {
		return err
	}

	a.Println(string(event))

	return nil
}
