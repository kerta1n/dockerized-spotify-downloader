cd /
go get github.com/go-flac/flacpicture github.com/go-flac/flacvorbis github.com/go-flac/go-flac github.com/valyala/fasthttp github.com/zmb3/spotify golang.org/x/oauth2/clientcredentials
go build -o metadata.cp -ldflags "-s -w" metadata.go
cp metadata.cp metadata
