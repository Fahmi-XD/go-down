FROM golang:1.23-alpine AS builder

WORKDIR /app

# RUN apk add --no-cache git

COPY . .

RUN go mod tidy && go build -o main .

FROM alpine:latest AS runtime

WORKDIR /app

RUN apk add --no-cache \
    chromium \
    chromium-chromedriver \
    ca-certificates \
    tzdata \
    dbus \
    dumb-init \
    font-noto \
    font-noto-cjk \
    font-noto-emoji \
    xvfb \
    xauth

ENV CHROME_BIN=/usr/bin/chromium
ENV CHROME_FLAGS="--headless --no-sandbox --disable-gpu --disable-software-rasterizer --disable-dev-shm-usage"
ENV XDG_CONFIG_HOME=/tmp/.chromium
ENV XDG_CACHE_HOME=/tmp/.chromium
ENV DISPLAY=:99

RUN mkdir -p /tmp/.chromium && chmod -R 777 /tmp/.chromium
RUN mkdir -p /storage/images

COPY --from=builder /app/main .
COPY --from=builder /app/database /app/database
COPY --from=builder /app/service /app/service
EXPOSE 8080

ENTRYPOINT ["/usr/bin/dumb-init", "--"]

CMD ["sh", "-c", "Xvfb :99 -screen 0 1920x1080x24 -shmem & ./main"]
# CMD ["sh", "-c", "Xvfb :99 -screen 0 1920x1080x24 & ./main"]