// Package slack15 is log15 (https://github.com/inconshreveable/log15)
// handler for sending log messages to Slack.
package slack15

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/inconshreveable/log15"
)

// ErrNoWebHook is returned if no WebHook URL is provided nor it
// could be found in environment
var ErrNoWebHook = errors.New("No Slack WebHook URL specified")

// Handler implements log15.Handler interface
type Handler struct {
	// WebHook URL (if empty taken from $SLACK_WEBHOOK)
	URL string

	// Message formatter – if nil default will be used
	Formatter log15.Format

	// Envelope The following fields allow for ovewritting default values
	// for webhook (as set in slack.com/services)
	Envelope
}

// func NewHandler()

type ctxReader struct {
	ctx []interface{}

	key   string
	value interface{}
	err   error
}

func (r *ctxReader) Pairs() int {
	return len(r.ctx) / 2
}

func (r *ctxReader) Next() bool {
	if len(r.ctx) < 2 {
		return false
	}
	var ok bool
	r.key, ok = r.ctx[0].(string)
	if !ok {
		r.err = fmt.Errorf("%+v is not a string key", r.ctx[0])
		r.key = "?"
	}
	r.value = r.ctx[1]
	r.ctx = r.ctx[2:]
	return true
}

func newCtxReader(ctx []interface{}) *ctxReader {
	return &ctxReader{ctx: ctx}
}

func (r *ctxReader) Key() string {
	return r.key
}

func (r *ctxReader) Value() interface{} {
	return r.value
}

func (r *ctxReader) Err() error {
	return r.err
}

// Log logs records by sending it to Slack
func (h *Handler) Log(r *log15.Record) error {
	msg, err := h.getMsg(r)
	// send message anyway if error occured

	url := h.URL
	if url == "" {
		url = os.Getenv("SLACK_WEBHOOK_URL")
		if url == "" {
			return ErrNoWebHook
		}
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("slack responsed with code %d", resp.StatusCode)
	}

	return err
}

func (h *Handler) getMsg(r *log15.Record) (*message, error) {
	var err error
	msg := &message{
		Envelope: h.Envelope,
	}

	color := "#32C8C8" // blue
	switch r.Lvl {
	case log15.LvlInfo:
		color = "good" // green
	case log15.LvlWarn:
		color = "warning" // yellow
	case log15.LvlError:
		color = "danger" // red
	case log15.LvlCrit:
		color = "#C832C8" // purple
	}

	if h.Formatter != nil {
		txt := string(h.Formatter.Format(r))
		msg.Attachments = []attachment{{
			Text:     txt,
			Fallback: txt,
			Color:    color,
		}}
	} else {
		ctx := newCtxReader(r.Ctx)
		fields := make([]field, 0, ctx.Pairs()+1)

		for ctx.Next() {
			v := fmt.Sprint(ctx.Value())
			fields = append(fields, field{
				Title: ctx.Key(),
				Value: v,
				Short: true,
			})
		}
		err = ctx.Err()

		msg.Attachments = []attachment{{
			Text:     r.Msg,
			Fallback: string(log15.LogfmtFormat().Format(r)),
			Fields:   fields,
			Color:    color,
		}}
	}

	return msg, err
}
