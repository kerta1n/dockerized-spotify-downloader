cd /
go build -o argvMatey.cp -ldflags "-s -w" argvMatey.go
cp argvMatey.cp argvMatey
