FROM  --platform=$BUILDPLATFORM golang:1.23-alpine AS builder
COPY . /go/src/github.com/xireiki/AhaDNS
WORKDIR /go/src/github.com/xireiki/AhaDNS
ARG TARGETOS TARGETARCH
ENV CGO_ENABLED=0
ENV GOOS=$TARGETOS
ENV GOARCH=$TARGETARCH
RUN set -ex \
    && apk add git build-base \
    && go build -o /go/bin/ahadns -v -trimpath
FROM  --platform=$BUILDPLATFORM alpine AS dist
RUN set -ex \
    && apk upgrade \
    && apk add bash tzdata ca-certificates \
    && rm -rf /var/cache/apk/*
COPY --from=builder /go/bin/ahadns /usr/local/bin/ahadns
ENTRYPOINT ["ahadns"]
