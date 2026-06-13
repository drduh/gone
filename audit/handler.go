package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

const errMarshal = `{"time":%q,` +
	`"level":"ERROR",` +
	`"message":"failed to marshal event",` +
	`"event":%q,` +
	`"error":%q}`

// Handle formats and prints audit events in JSON format.
func (a *Auditor) Handle(
	_ context.Context,
	r slog.Record,
) error {
	data := make(map[string]any, r.NumAttrs())
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
		a.Printf(errMarshal,
			r.Time.Format(a.TimeFormat),
			r.Message,
			err.Error(),
		)
		return fmt.Errorf("marshal log record: %w", err)
	}

	a.Println(string(event))

	return nil
}
