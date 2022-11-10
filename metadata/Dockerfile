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

ENV PLAYLIST_ID='TEST_PLAYLIST_ID'  
ENV SPOTIFY_ID='TEST_SPOTIFY_ID'
ENV SPOTIFY_SECRET='TEST_SECRET'

CMD ["./metadata"]