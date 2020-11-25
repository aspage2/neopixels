FROM balenalib/raspberrypi3-golang:latest-build AS builder
RUN [ "cross-build-start" ]
WORKDIR /tmp
RUN apt-get update -y && apt-get install -y scons
RUN git clone --depth=1 https://github.com/jgarff/rpi_ws281x.git && \
  cd rpi_ws281x && \
  scons
RUN [ "cross-build-end" ]

FROM balenalib/raspberrypi3-golang:latest
RUN [ "cross-build-start" ]
COPY --from=builder /tmp/rpi_ws281x/*.a /usr/local/lib/
COPY --from=builder /tmp/rpi_ws281x/*.h /usr/local/include/
RUN go get -v -u github.com/rpi-ws281x/rpi-ws281x-go
RUN go get -v -u google.golang.org/grpc
RUN go get -v -u github.com/spf13/viper
RUN [ "cross-build-end" ]
