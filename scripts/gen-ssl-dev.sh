#!/bin/sh
# Generate self-signed SSL cert for local development only.
# For production use Let's Encrypt / Certbot.
set -e

OUT=nginx/ssl
mkdir -p "$OUT"

openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout "$OUT/privkey.pem" \
  -out "$OUT/fullchain.pem" \
  -subj "/C=NG/ST=Lagos/L=Lagos/O=SiteLogix Dev/CN=localhost"

echo "Self-signed cert written to $OUT/"
echo "WARNING: This is for local development only. Use a real CA cert in production."
