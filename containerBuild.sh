go mod init argvMatey
go mod tidy
go build -o argvMatey.cp -ldflags "-s -w" argvMatey.go
cp argvMatey.cp argvMatey
rm argvMatey.cp