# https://hub.docker.com/_/golang
FROM docker.io/library/golang:1.25.3-alpine AS build

ARG GOOS=linux

ARG GOARCH=amd64

ENV GOOS=${GOOS}

ENV GOARCH=${GOARCH}

WORKDIR /src

COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o webhook -ldflags "-w -s -extldflags '-static' ${versionflags}" .

FROM gcr.io/distroless/static-debian12:nonroot AS webhook

USER nonroot

COPY --from=build /src/webhook /bin/webhook

ENTRYPOINT ["/bin/webhook"]
