FROM alpine:latest as runner
COPY output/app /usr/local/bin/app
CMD ["app"]
