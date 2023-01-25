
### Build ###
FROM golang:1.19.5-alpine3.17 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /comicbot

### Deploy

FROM alpine

WORKDIR /

COPY --from=build /comicbot /comicbot

ENTRYPOINT ["/comicbot"]