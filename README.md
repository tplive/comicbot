# Lunchbot!

Go program that goes out to the [Teknisk Ukeblad web](https://tu.no) and downloads todays Lunch cartoon, then sends it to your Slack channel. The cartoons are in Norwegian.

TU also publishes Dilbert cartoons in parallel. To get Dilbert just replace the value for `comicId` in the code with `dilbert`.

The original idea is mine, but I have used various resources in implementing the solution. Learning Go along the way.

You should probably avoid spamming the tu.no website, as they might react with breaking changes.

## Prerequisites
1. Install Golang
1. Set environment variable as below, or in a .env file for automatic inclusion:
```bash
WEBHOOK_URL="https://hooks.slack.com/services/THIS/IS/PRIVATE"
```
3. Run the bot with `go run lunchbot.go`

This will download todays cartoon, if there is one, or fail gracefully otherwise.


