FROM golang:1.16-buster AS build
WORKDIR /src
COPY ./ ./
RUN go build go-nest-temp-monitor.go

FROM debian:buster
RUN apt-get update && apt-get install -y ca-certificates
COPY --from=build /src/go-nest-temp-monitor /

CMD [ "/go-nest-temp-monitor", "-c", "/app/config.json" ]
