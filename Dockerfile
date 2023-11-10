FROM alpine
RUN apk add --no-cache ca-certificates curl
COPY clouddriver /usr/local/bin
RUN ls -l /usr/local/bin
CMD ["/usr/local/bin/clouddriver"]
