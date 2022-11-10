FROM ubuntu
RUN apt update && \
apt install alsa-utils curl ffmpeg pulseaudio -y && \
curl -L https://github.com/Spotifyd/spotifyd/releases/download/v0.2.23/spotifyd-linux-slim.tar.gz -o spotifyd.tar.gz && \
tar xzf spotifyd.tar.gz && rm spotifyd.tar.gz
COPY entrypoint.sh onChange spotifyd.conf argvMatey /
CMD ["./entrypoint.sh"]