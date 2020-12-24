setup:
	go get .
	go build .
	cp .env.example .env

build:
	go get .
	go build .

compile:
	echo "Compiling for most os's and platforms"
	echo "-> linux:"
	GOOS=linux GOARCH=386 go build -o bin/dcpresencebot-linux-386 main.go
	GOOS=linux GOARCH=amd64 go build -o bin/dcpresencebot-linux-amd64 main.go
	GOOS=linux GOARCH=arm go build -o bin/dcpresencebot-linux-arm main.go
	GOOS=linux GOARCH=arm64 go build -o bin/dcpresencebot-linux-arm64 main.go
	echo "-> darwin:"
	GOOS=darwin GOARCH=386 go build -o bin/dcpresencebot-darwin-386 main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/dcpresencebot-darwin-amd64 main.go
	# GOOS=darwin GOARCH=arm go build -o bin/dcpresencebot-darwin-arm main.go
	# GOOS=darwin GOARCH=arm64 go build -o bin/dcpresencebot-darwin-arm64 main.go
	echo "-> freebsd:"
	GOOS=freebsd GOARCH=arm go build -o bin/dcpresencebot-freebsd-arm main.go
	GOOS=freebsd GOARCH=amd64 go build -o bin/dcpresencebot-freebsd-amd64 main.go
	GOOS=freebsd GOARCH=386 go build -o bin/dcpresencebot-freebsd-386 main.go
	echo "-> windows:"
	GOOS=windows GOARCH=amd64 go build -o bin/dcpresencebot-windows-amd64.exe main.go
	GOOS=windows GOARCH=386 go build -o bin/dcpresencebot-windows-i386.exe main.go

linuxcrossmac:
	GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc CGO_ENABLED=0 go build -a -v -ldflags '-d -s -w' -installsuffix cgo -o bin/dcpresencebot-linux-amd64-cgo main.go

run:
	go run .
