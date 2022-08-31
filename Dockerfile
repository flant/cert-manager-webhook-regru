FROM golang:1.18.3-alpine AS build_deps

WORKDIR /src

RUN apk add --no-cache git
COPY go.* .

RUN go mod download


FROM build_deps AS build

COPY . .

RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -extldflags "-static"' .

FROM alpine:3.9

RUN apk add --no-cache ca-certificates

COPY --from=build /workspace/webhook /usr/local/bin/webhook

ENTRYPOINT ["webhook"]