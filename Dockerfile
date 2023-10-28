# Build app
FROM golang:1-alpine3.17 AS go-app

WORKDIR /app

COPY . /app
RUN go mod vendor \
    && go build -o templar

# =========
FROM alpine:3.17

LABEL author="bravepickle"
LABEL version="1.0"
LABEL description="CLI templatator"

WORKDIR /app

COPY --from=go-app /app/templar /app/templar

VOLUME /app

ENTRYPOINT [ "/app/templar" ]
