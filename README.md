# Comicbot!

Go program that goes out to the [Teknisk Ukeblad web](https://tu.no) and downloads todays Lunch and Dilbert cartoons, then sends the URL to your Slack channel. The cartoons are in Norwegian.

Could evolve to get other comics as well? :)

The original idea is mine, but I have used various resources in implementing the solution. Learning Go along the way.

You should probably avoid spamming the tu.no website, as they might react with breaking changes.

## Prerequisites

1. Get your Slack Incoming Webhook URL set up.
1. Install Golang
1. Set environment variable as below, or in a .env file for automatic inclusion:

```bash
WEBHOOK_URL="https://hooks.slack.com/services/THIS/IS/PRIVATE"
```

3. Run the bot with `go run lunchbot.go`

This will download todays cartoon, if there is one, or fail gracefully otherwise.

## Building and running with Docker

1. Build the image `docker build --tag comicbot .`
1. Run the container `docker run --env WEBHOOK_URL="<your-slack-webhook-url>" comicbot`
