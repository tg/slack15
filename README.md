# slack15 – from log15 to Slack [![GoDoc](https://godoc.org/github.com/tg/slack15?status.svg)](https://godoc.org/github.com/tg/slack15)

Say hello to Slack handler for [log15](https://github.com/inconshreveable/log15) and
send log messages in Go directly to Slack via [incoming webhook](https://api.slack.com/incoming-webhooks).

You will need at least webhook URL to make it working.

## Quick start
Pass webhook URL to `slack15.Handler` or set it up in `SLACK_WEBHOOK_URL` environmental variable.

```go
log := log15.New()
log.SetHandler(&slack15.Handler{
	URL: "", // pass it here, or not
})
log.Info("Whaam!", "who", "Roy Lichtenstein", "when", 1963)
```

![screenshot](https://www.dropbox.com/s/0v27ia73ox100gd/github_slack15.png?raw=1)
