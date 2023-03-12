FROM golang:alpine
RUN  apk --no-cache add build-base mosquitto-dev

WORKDIR /src
ARG CGO_ENABLED=1
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
#COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v .

FROM alpine:latest
RUN apk --no-cache add mosquitto
WORKDIR /root/
COPY --from=0 /src/sensor ./sensor

ENTRYPOINT [ "./sensor" ]
CMD ["-user", "mqtt", "-password", "mqtt", "-host", "http://host.docker.internal:1883"]