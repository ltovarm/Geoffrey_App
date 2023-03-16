FROM golang:alpine
RUN  apk --no-cache add build-base

WORKDIR /src
ARG CGO_ENABLED=1
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v .

FROM alpine:latest

WORKDIR /root/
COPY --from=0 /src/consumer ./consumer

ENTRYPOINT [ "./consumer" ]
CMD ["-broker", "amqp://mqtt:mqtt@host.docker.internal"]