FROM docker.io/library/golang:1.23-alpine AS build_deps
WORKDIR /src
COPY . .
RUN go mod download
RUN go get -u golang.org/x/crypto@v0.31.0 # Fix https://avd.aquasec.com/nvd/cve-2024-45337
RUN go get -u golang.org/x/net@0.33.0 # https://avd.aquasec.com/nvd/cve-2024-45338
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o webhook -ldflags '-w -s -extldflags "-static"' .

FROM gcr.io/distroless/static-debian12:latest
COPY --from=build_deps /src/webhook /bin/webhook
USER nonroot
ENTRYPOINT ["/bin/webhook"]
