package slack15_test

import (
	"github.com/inconshreveable/log15"
	"github.com/tg/slack15"
)

func ExampleHandler() {
	log := log15.New()
	log.SetHandler(&slack15.Handler{
		URL: "", // pass url here, or through $SLACK_WEBHOOK_URL
		// You can skip this and stick with webhook defaults
		Envelope: slack15.Envelope{
			Username:  "Mr. Bellamy",
			IconEmoji: ":showman:",
		},
	})
	log.Info("Whaam!", "who", "Roy Lichtenstein", "when", 1963)

	// Output:
}
