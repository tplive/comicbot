# Comicbot!

Go program that goes out to the [Teknisk Ukeblad web](https://tu.no) and downloads todays Lunch and Dunce cartoons, and Dilbert - while supplies last[1] - then sends the URL to your Slack channel. The cartoons are in Norwegian.

New! Also supports downloading XKCD - will use a key/value pair at [kvdb.io](https://kvdb.io/) to track the last comic that was downloaded, and only download newer ones.

The original idea is mine, but I have used various resources in implementing the solution. Learning Go along the way.

You should probably avoid spamming the tu.no website, as they might react with breaking changes.

[1] Dilbert was cancelled following the Scott Adams controversy.

## Prerequisites

1. Get your Slack Incoming Webhook URL set up.
2. Set up your kvdb.io bucket, ie like this `curl -d 'email=user@example.com' https://kvdb.io`
3. Install Golang. Or just run the bot in a Github workflow. See `.github/workflows` for examples.
4. Set environment variables as below, or in a `.env` file for automatic inclusion for local development:

   ```bash
   WEBHOOK_URL="https://hooks.slack.com/services/YOUR/PRIVATE/PARTS"
   KVDB_BUCKET="yourBucketId"
   ```

5. Set a key in your KVDB bucket to the comic id of the XKCD comic you want to start tracking from. As of this writing the current comic is 2752.

   ```bash
   curl https://kvdb.io/yourBucketId/xkcd -d '2752'
   ```

   Attention!! If you set this number to anything below the latest comic, it will attempt to download every comic sequentially!

6. Run the bot with `go run .`

This should download the images of the current comics to `pwd`, and post their URLs to Slack.

## Building and running with Docker

1. Build the image `docker build --tag comicbot .`
1. Run the container `docker run --env WEBHOOK_URL="<your-slack-webhook-url>" --env KVDB_BUCKET="<yourBucketId>" comicbot`

## Building and publishing to Dockerhub

First login to Dockerhub with `docker login`

1. Build the image and tag appropriately `docker build --tag <repository>/comicbot:latest .`
1. Push the image to the repository `docker push <repository>/comicbot:latest`
