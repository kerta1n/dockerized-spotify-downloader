package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

var extension = ".flac"
var maxSongLength = time.Hour
var volumePath = "/volume/"

func endContainer(err error, l *log.Logger) {
	l.Println("time to die")
	l.Println(err)
	out, _ := exec.Command("pgrep", "spotifyd").Output()
	pid, _ := strconv.Atoi(strings.TrimSpace(string(out)))
	syscall.Kill(pid, syscall.SIGKILL)
	panic(err)
}

func getEventOrTrack(fileName string, l *log.Logger) string {
	s, err := ioutil.ReadFile(fileName)
	if err != nil {
		endContainer(errors.New(fmt.Sprintf("could not read from file: %s", fileName)), l)
	}
	return strings.TrimSpace(string(s))
}

func main() {
	l := log.New(&bytes.Buffer{}, "argvMatey: ", log.LUTC)
	l.SetOutput(os.Stdout)
	l.Println("starting argvMatey")
	sigUsr := make(chan os.Signal, 1)
	signal.Notify(sigUsr, syscall.SIGUSR1)
	death := make(chan struct{})
	for {
		select {
		case <-sigUsr:
			close(death)
			death = make(chan struct{})
			event := getEventOrTrack("PLAYER_EVENT", l)
			track := getEventOrTrack("TRACK_ID", l)
			l.Println(fmt.Sprintf("Got signal. Event: %s Track: %s", event, track))
			switch event {
			case "start":
				go record(death, l, track)
			case "change":
				go record(death, l, track)
			case "stop":
				endContainer(errors.New("spotify stopped"), l)
				return
			}
		}
	}
}

func record(death chan struct{}, l *log.Logger, track string) {
	l.Println(fmt.Sprintf("track: %s", track))
	fileName := volumePath + track + extension
	cmd := exec.Command("ffmpeg", "-y", "-f", "pulse", "-i", "default", fileName)
	err := cmd.Start()
	if err != nil {
		l.Println("failed to record track")
		return
	}
	l.Println("recording started")
	wgC := make(chan struct{})
	go func() {
		defer close(wgC)
		time.Sleep(maxSongLength)
	}()
	select {
	case <-wgC:
		endContainer(errors.New("song is either longer than one hour or this broke"), l)
	case <-death:
		break
	}
	err = cmd.Process.Signal(os.Interrupt)
	if err != nil {
		l.Println("failed to stop recording.")
		endContainer(err, l)
	}
	l.Println("recording stopped")
}
