FROM golang:1.18.3-buster AS build_deps
WORKDIR /src
RUN DEBIAN_FRONTEND=noninteractive; apt-get update && apt-get install git -y
ENV GO111MODULE=on
COPY . .
RUN ls -la /src
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .


FROM debian:buster-slim
RUN DEBIAN_FRONTEND=noninteractive; apt-get update && apt-get install ca-certificates -y
COPY --from=build_deps /src/webhook /usr/local/bin/webhook
ENTRYPOINT ["webhook"]