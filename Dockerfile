
### Build ###
FROM golang:1.17.1-alpine3.14 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /comicbot

### Deploy

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=build /comicbot /comicbot

USER nonroot:nonroot

ENTRYPOINT ["/comicbot"]