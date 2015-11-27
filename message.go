package slack15

// Envelope describes destination and sender.
// Defaults are defined by webhook setting in Slack and this structure
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

// Msg represents Slack message
type Msg struct {
	Envelope
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

// Attachment represents attachment to a message
type Attachment struct {
	Fallback string  `json:"fallback,omitempty"` // Required plain-text summary of the attachment.
	Color    string  `json:"color,omitempty"`
	Pretext  string  `json:"pretext,omitempty"`
	Text     string  `json:"text,omitempty"`
	Fields   []Field `json:"fields,omitempty"`
}

// Field represents field withing an attachment
type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}
