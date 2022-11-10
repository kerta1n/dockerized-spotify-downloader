FROM golang:1.19.2 AS builder

ENV APP_HOME /go/src/metadata

WORKDIR "$APP_HOME"
COPY . .
RUN go mod init
RUN go mod tidy
RUN go mod download metadata
RUN go mod verify
RUN go build -o metadata

FROM golang:1.19.2

ENV APP_HOME /go/src/metadata
RUN mkdir -p "$APP_HOME"
WORKDIR "$APP_HOME"

COPY --from=builder "$APP_HOME"/metadata $APP_HOME

ARG PLAYLIST_ID='CHANGEME'
ARG SPOTIFY_ID='CHANGEME'
ARG SPOTIFY_SECRET='CHANGEME'

CMD ["./metadata"]
