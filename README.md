# Lunchbot!

Slackbot that goes out to the [Teknisk Ukeblad web](https://tu.no) and downloads todays Lunch cartoon, then sends it to your Slack channel. The cartoons are in Norwegian.

TU also publishes Dilbert cartoons in parallel. To get Dilbert just replace the value for `comicId` in the code with `dilbert`.

The original idea is mine, but I have used various resources in implementing the solution. Learning Go along the way.

## Pre-requisites
1. Install Golang
1. Run the bot with `go run lunchbot.go`

This will download todays cartoon, if there is one, or fail otherwise.
