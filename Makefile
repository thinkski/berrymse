all: armv6l armv7l

armv6l: pkged.go
	mkdir -p armv6l
	GOARCH=arm GOARM=6 GOOS=linux go build -v -o armv6l/berrymse -ldflags="-w -s"

armv7l: pkged.go
	mkdir -p armv7l
	GOARCH=arm GOARM=7 GOOS=linux go build -v -o armv7l/berrymse -ldflags="-w -s"

pkged.go:
	pkger

clean:
	rm -rf armv6l armv7l pkged.go

.PHONY: armv6l armv7l
.PHONY: clean
