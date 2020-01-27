touch argvMatey
docker run --rm -it -v ${PWD}/containerBuild.sh:/containerBuild.sh -v ${PWD}/argvMatey.go:/argvMatey.go -v ${PWD}/argvMatey:/argvMatey golang bash /containerBuild.sh
chmod +x argvMatey
docker network create --subnet=172.18.0.0/16 workwithvpn
