// Package slack15 is log15 (https://github.com/inconshreveable/log15)
// handler for sending log messages to Slack.
package slack15

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/inconshreveable/log15.v2"
)

// ErrNoWebHook is returned if no WebHook URL is provided nor it
// could be found in environment
var ErrNoWebHook = errors.New("No Slack WebHook URL specified")

// Handler implements log15.Handler interface
type Handler struct {
	// Envelope allows to overwrite webhook's defaults
	Envelope

	// WebHook URL (if empty taken from $SLACK_WEBHOOK)
	URL string

	// Message formatter – if nil default will be used
	Formatter log15.Format

	// Client is used for HTTP requests to Slack API.
	// If nil, http.DefaultClient is used.
	Client *http.Client
}

// Log logs records by sending it to Slack
func (h *Handler) Log(r *log15.Record) error {
	msg := h.getMsg(r)

	// Take URL from handler or environment
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

	c := h.Client
	if c == nil {
		c = http.DefaultClient
	}

	resp, err := c.Post(url, "", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Report up to 1kB of response body
		body, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 1024))
		err = fmt.Errorf("slack responded with %d: %s", resp.StatusCode, body)
	}

	return err
}

// getMsg returns message which should be sent to Slack
func (h *Handler) getMsg(r *log15.Record) *message {
	msg := &message{
		Envelope: h.Envelope,
	}

	// Choose message color depending on log level
	// (this is imitating colors from log15.TerminalFormat)
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
	msg.Attachments[0].Color = color

	if h.Formatter != nil {
		txt := string(h.Formatter.Format(r))
		msg.Attachments[0].Text = txt
		msg.Attachments[0].Fallback = txt
	} else {
		ctx := newCtxReader(r.Ctx)
		fields := make([]Field, 0, len(r.Ctx)/2)

		for ctx.Next() {
			v := ctx.Value()
			// See if value is a Field; if not, fill with defaults
			f, ok := v.(Field)
			if !ok {
				f.Value = fmt.Sprintf("%+v", v)
				f.Short = true
			}
			if f.Title == "" {
				f.Title = ctx.Key()
			}
			fields = append(fields, f)
		}

		msg.Attachments[0].Text = r.Msg
		msg.Attachments[0].Fallback = string(log15.LogfmtFormat().Format(r))
		msg.Attachments[0].Fields = fields
	}

	return msg
}
