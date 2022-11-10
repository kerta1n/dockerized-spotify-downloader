package main

import (
        "bytes"
        "context"
        "log"
        "os"
        "path/filepath"
        "strconv"
        "strings"
        "sync"

        "github.com/go-flac/flacpicture"
        "github.com/go-flac/flacvorbis"
        "github.com/go-flac/go-flac"
        "github.com/valyala/fasthttp"
        "github.com/zmb3/spotify"
        "golang.org/x/oauth2/clientcredentials"
)

var extension = ".flac"
var volume = "volume/"

type albumTrack struct {
        album spotify.SimpleAlbum
        track spotify.SimpleTrack
}

type dupCount struct {
        at    []*albumTrack
        count uint
}

type getImage struct {
        image *[]byte
}

func image(b *[]byte, fileName string) *flac.File {
        f, err := flac.ParseFile(fileName)
        if err != nil {
                panic(err)
        }
        picture, err := flacpicture.NewFromImageData(flacpicture.PictureTypeFrontCover, "Front cover", *b, "image/jpeg")
        if err != nil {
                panic(err)
        }
        picturemeta := picture.Marshal()
        f.Meta = append(f.Meta, &picturemeta)
        return f
}

func info(clientID, clientSecret string, l *log.Logger, playlistID spotify.ID) (map[string]*albumTrack, error) {
        config := &clientcredentials.Config{
                ClientID:     clientID,
                ClientSecret: clientSecret,
                TokenURL:     spotify.TokenURL,
        }
        token, err := config.Token(context.Background())
        if err != nil {
                return nil, err
        }
        client := spotify.Authenticator{}.NewClient(token)
        tracks, err := client.GetPlaylistTracks(playlistID)
        if err != nil {
                return nil, err
        }
        m := make(map[string]*albumTrack)
        for _, track := range tracks.Tracks {
                m[track.Track.SimpleTrack.ID.String()] = &albumTrack{
                        album: track.Track.Album,
                        track: track.Track.SimpleTrack,
                }
        }
        l.Printf("playlist has %d tracks", len(tracks.Tracks))
        duplicates := make(map[string]*dupCount, 0)
        for _, at := range m {
                name := at.track.Name
                if _, ok := duplicates[name]; !ok {
                        duplicates[name] = &dupCount{
                                at:    make([]*albumTrack, 0),
                                count: 0,
                        }
                }
                duplicates[name].at = append(duplicates[name].at, at)
                duplicates[name].count++
        }
        for k, v := range duplicates {
                if v.count > 1 {
                        l.Printf("Song \"%s\" appears %d times.", k, v.count)
                        sameId := make(map[string]bool)
                        for _, at := range v.at {
                                id := at.track.ID.String()
                                if _, ok := sameId[id]; !ok {
                                        sameId[id] = true
                                        artists := ""
                                        for _, artist := range at.track.Artists {
                                                artists += ", " + artist.Name
                                        }
                                        artists = strings.Trim(artists, ", ")
                                        artist := " (" + artists + ")"
                                        newName := k + artist
                                        at.track.Name = k + artist
                                        l.Printf("One \"%s\" will be renamed \"%s\"", k, newName)
                                } else {
                                        l.Printf("Same ID found more than once: %s", id)
                                }
                        }
                }
        }
        return m, nil
}

func main() {
        l := log.New(&bytes.Buffer{}, "metadata: ", log.LUTC)
        l.SetOutput(os.Stdout)
        l.Println("start metadata")
        clientID := os.Getenv("SPOTIFY_ID")
        clientSecret := os.Getenv("SPOTIFY_SECRET")
        playlistID := spotify.ID(os.Getenv("PLAYLIST_ID"))
        if len(clientID) == 0 || len(clientSecret) == 0 || len(playlistID) == 0 {
                l.Println("Not all environment variables set.")
                return
        }
        files, err := filepath.Glob(volume + "*" + extension)
        if err != nil {
                panic(err)
        }
        if len(files) == 0 {
                l.Println("No files in: " + volume)
                return
        }
        l.Printf("have %d files", len(files))
        m, err := info(clientID, clientSecret, l, playlistID)
        if err != nil {
                panic(err)
        }
        wg := &sync.WaitGroup{}
        for _, fileName := range files {
                id := strings.Split(strings.Split(fileName, ".")[0], "/")[1]
                at, ok := m[id]
                if !ok {
                        l.Printf("Couldn't find info for %s", fileName)
                        continue
                }
                newName := volume + strings.ReplaceAll(at.track.Name, "/", "") + extension
                wg.Add(1)
                go meta(at, fileName, l, newName, wg)
        }
        wg.Wait()
}

func meta(at *albumTrack, fileName string, l *log.Logger, newName string, wg *sync.WaitGroup) {
        defer wg.Done()
        wg2 := &sync.WaitGroup{}
        gi := &getImage{}
        if len(at.album.Images) > 0 {
                l.Println("downloading album image: " + newName)
                wg2.Add(1)
                go func() {
                        _, imgData, _ := fasthttp.Get(make([]byte, 0), at.album.Images[0].URL)
                        gi.image = &imgData
                        wg2.Done()
                }()
        }
        err := os.Rename(fileName, newName)
        if err != nil {
                panic(err)
        }
        var f *flac.File
        if len(at.album.Images) > 0 {
                wg2.Wait()
                l.Println("inserting image: " + newName)
                f = image(gi.image, newName)
        }
        l.Println("inserting vorbis: " + newName)
        vorbis(at, f, newName)
        l.Println("done: " + newName)
}

func vorbis(at *albumTrack, f *flac.File, fileName string) {
        var err error
        if f == nil {
                f, err = flac.ParseFile(fileName)
                if err != nil {
                        panic(err)
                }
        }
        cmts := flacvorbis.New()
        err = cmts.Add(flacvorbis.FIELD_ALBUM, at.album.Name)
        if len(at.track.Artists) > 0 {
                artists := ""
                for _, artist := range at.track.Artists {
                        artists += ", " + artist.Name
                }
                artists = strings.Trim(artists, ", ")
                err = cmts.Add(flacvorbis.FIELD_ARTIST, artists)
        }
        err = cmts.Add(flacvorbis.FIELD_DATE, at.album.ReleaseDate)
        url, ok := at.track.ExternalURLs["spotify"]
        if ok {
                err = cmts.Add(flacvorbis.FIELD_DESCRIPTION, url)
        }
        err = cmts.Add(flacvorbis.FIELD_TRACKNUMBER, strconv.Itoa(at.track.TrackNumber))
        err = cmts.Add(flacvorbis.FIELD_TITLE, at.track.Name)
        if err != nil {
                panic(err)
        }
        cmtsmeta := cmts.Marshal()
        f.Meta = append(f.Meta, &cmtsmeta)
        err = f.Save(fileName)
        if err != nil {
                panic(err)
        }
}