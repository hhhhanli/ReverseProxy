go mod tidy
export GOOS=linux
export GOARCH=amd64
go build main.go
zip main.zip main
