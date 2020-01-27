touch metadata
docker run --rm -it -v ${PWD}/containerBuild.sh:/containerBuild.sh -v ${PWD}/metadata.go:/metadata.go -v ${PWD}/metadata:/metadata golang bash /containerBuild.sh
chmod +x metadata
