@rem go build testsock1.go -ldflags="--subsystem windows"
go build -ldflags "-H windowsgui -s -w" %1