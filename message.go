package slack15

import "fmt"

// Envelope describes destination and sender (as visiable in Slack).
// Defaults are defined by webhook settings in Slack and this structure
// allows for overwritting these.
// More info at https://api.slack.com/incoming-webhooks
type Envelope struct {
	// Destination channel in slack, e.g. "#log" or "@bobby"
	Channel string `json:"channel,omitempty"`

	// Username
	Username string `json:"username,omitempty"`

	// Icon URL, e.g. "https://slack.com/img/icons/app-57.png"
	IconURL string `json:"icon_url,omitempty"`

	// Icon emoji, e.g. ":rabbit:"
	IconEmoji string `json:"icon_emoji,omitempty"`
}

// Field represents a key-value pair within slack message
type Field struct {
	Title string `json:"title"`           // Field title (key)
	Value string `json:"value"`           // Text
	Short bool   `json:"short,omitempty"` // If false field will occupy the whole row
}

// String returns field's value, so it's correctly printed in other formatters
func (f Field) String() string {
	return f.Value
}

// Long returns Field which will occupy the whole row in Slack message.
// It can be used as a value of key-value pair passed to logger.
//
// Example:
// log.Info("got message", "source", ip, "data", slack15.Long(data)).
func Long(value interface{}) Field {
	return Field{
		Title: "", // will use key from log context
		Value: fmt.Sprint(value),
		Short: false,
	}
}

// Message represents a single slack message.
// It's serialised to JSON and send to incoming webhook.
type message struct {
	Envelope
	Attachments [1]struct {
		Text     string  `json:"text,omitempty"`
		Fallback string  `json:"fallback,omitempty"` // Required plain-text summary
		Color    string  `json:"color,omitempty"`
		Fields   []Field `json:"fields,omitempty"`
	} `json:"attachments,omitempty"`
}
