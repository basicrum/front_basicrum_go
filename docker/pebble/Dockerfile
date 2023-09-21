FROM letsencrypt/pebble:latest
RUN apk update
RUN apk add --update bash curl
COPY check.sh check.sh
RUN chmod +x check.sh
