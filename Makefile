
bin/%: cmd/%/main.go $(wildcard internal/*.go) $(wildcard internal/pattern/*.go) ledproto/leds.pb.go
	@mkdir -p bin
	docker run --rm -it \
	-v "$(shell pwd)":/go/src/pnpleds \
	rpi-leds \
	/usr/bin/qemu-arm-static /bin/sh -c "go build -o src/pnpleds/bin/$(@F) -v pnpleds/cmd/$(@F)"

upload: $(wildcard bin/*)
	scp bin/* pi@192.168.2.25:~

proto: ledproto/leds.pb.go
.PHONY: proto

ledproto/leds.pb.go: leds.proto
	@mkdir -p ledproto
	protoc -I. --go_out=ledproto --go-grpc_out=ledproto leds.proto

fmt:
	goimports -w .

clean:
	rm -rf internal/ledproto/
	rm -rf bin/
	rm -f scripts/*pb2*.py

.PHONY: clean
