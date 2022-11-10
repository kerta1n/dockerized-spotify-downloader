pulseaudio -D --exit-idle-time=-1
pacmd load-module module-virtual-sink sink_name=v1
pacmd set-default-sink v1
pacmd set-default-source v1.monitor
cp /spotifyd.conf /spotifyd.conf.using
echo "device_name = $DEVICE_NAME" >>/spotifyd.conf.using
echo "password = $PASSWORD" >>/spotifyd.conf.using
echo "username = $USERNAME" >>/spotifyd.conf.using
./argvMatey &
./spotifyd --no-daemon --config-path /spotifyd.conf.using
