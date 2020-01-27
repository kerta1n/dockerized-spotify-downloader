# dockerized-spotify-downloader
Here is a self hosted solution to download Spotify playlists in a Docker container silently spinning in the background.
This method requires a user to use Spotify to switch playback to this downloader's device and start the playlist. After
that you must wait for the playlist to finish. Optionally, you can grab the metadata from Spotify and automatically
insert it into the audio files. The bitrate is customizable, it defaults to 320. This project should work as long as
[spotifyd](https://github.com/Spotifyd/spotifyd) works.

Disclaimer: this is a guide only for technical people.

## Prerequisites
* A Linux machine that runs: [Docker](https://www.docker.com/) with sufficient disk space.
* An internet connection good enough to stream audio from Spotify.
* A Spotify premium account. (Free trial instructions below.)
* Optional: Spotify API `Client ID` and `Client Secret` (for metadata).

## Quick start
**1** Edit the `start.sh` and `metadata/start.sh` files to set the environment variables.

**2** Run these commands to start the Spotify device:
```bash
user@host:~/dockerized-spotify-downloader ./build.sh
```
If this kicks back `Error response from daemon: Pool overlaps with other one on this address space`, then you likely
either have a VPN that should be temporarily disabled to create a new Docker network *OR* you need to choose another IP
block because there is an overlapping Docker network.
```bash
user@host:~/dockerized-spotify-downloader ./start.sh
```
You'll need to use `docker-compose down` to clean up before using the container again because the way it shuts down is
messy.

**3** Make sure repeat is not on and start the playlist on the device whose name you set in `setup.sh`. 

**4** Wait for the playlist to complete. If you started it in a web browser, you can close that now. (Other starting methods not tested.) The container will exit on playback stop.

**5** (Optional) Change the owner of the audio files if your user isn't a member of a group allowed to use Docker.

**6** (Optional) Run these commands to put album art and other metadata in the audio files:
```bash
user@host:~/dockerized-spotify-downloader cd metadata && ./build.sh && cd ..
user@host:~/dockerized-spotify-downloader/metadata$ ./metadata/start.sh
```
The working dir needs to contain the `volume` dir with the audio files in it.

**7** Grab the `.flac` files from the `~/dockerized-spotify-downloader/volume` directory. If you don't like `.flac` files, use [ffmpeg](https://www.ffmpeg.org/) to convert them.

## Environment variables
`start.sh`:
* `export DEVICE_NAME=` The name you want the Spotify device to be.
* `export PASSWORD=` Your Spotify premium account password.
* `export USERNAME=` Your Spotify premium account username. This is different from your display name. Find it in your `Account overview`.

`metadata/start.sh`
* `export PLAYLIST_ID=` The playlist ID. You can find this in the URL link to the playlist.
* `export SPOTIFY_ID=` Your Spotify developer `Client ID`.
* `export SPOTIFY_SECRET=` Your Spotify developer `Client Secret`.

## Getting a free Spotify premium trial
Clear all your cookies or use a private browser. Sign up for a Spotify account using a new disassociated email. After
this, there should be an option to have 30 days of free premium membership. Sign up for this, when they ask for a credit
card, use a service that creates a temporary credit card like [privacy.com](https://privacy.com) and allocate $1. Save
that username and password. Validate the email. Create a Spotify developer account and application under this account.
Remember to cancel Spotify premium before the end of the 30 trial.

___

## How to improve this project
* Remove user interaction by passing the playlist to the program and have the playlist auto start on container startup. See [this GitHub issue](https://github.com/Spotifyd/spotifyd/issues/78) or the Spotify Web API.
* Record the audio from `spotifyd` directly by using `--pipe` instead of a pulseaudio loopback. See [this GitHub issue](https://github.com/Spotifyd/spotifyd/issues/78).
* Implement a cleaner way to stop the container.

___

## Privacy
Everything in this project is open source and can be manually verified. However, to stay private, make sure to use a
reputable paid VPN service, follow best web browser privacy practices, and don't include have anything that can link
back to you. If you want a playlist, it may be best to recreate the playlist in the new disassociated account, with some
different songs and split into multiple playlists. This project is for private use only. Be wary of sharing any
downloaded songs, it's trivial for Spotify to add near invisible and unique trackers to every song determine where any
song came from.
