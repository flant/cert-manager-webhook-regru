# https://hub.docker.com/_/golang
FROM docker.io/library/golang:1.23.4-alpine AS build_deps
ARG GOOS=linux
ARG GOARCH=amd64
ENV GOOS=${GOOS}
ENV GOARCH=${GOARCH}
WORKDIR /src
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o webhook -ldflags '-w -s -extldflags "-static"' .

# https://github.com/GoogleContainerTools/distroless
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build_deps /src/webhook /bin/webhook
USER nonroot
ENTRYPOINT ["/bin/webhook"]
