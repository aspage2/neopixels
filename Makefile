build: bin/ledserver
.PHONY: build

bin/ledserver: $(wildcard ledserver/internal/*.go) $(wildcard ledserver/*.go)
	@mkdir -p bin
	docker run --rm --platform linux/arm/v7 \
	-v "$(shell pwd)":/go/src/pnpleds \
	rpi-leds \
	/usr/bin/qemu-arm-static /bin/sh -c "cd /go/src/pnpleds && go build -o bin/ledserver -v ./ledserver"

fmt:
	goimports -w .

clean:
	rm -rf internal/ledproto/
	rm -rf bin/
	rm -f scripts/*pb2*.py

.PHONY: clean
