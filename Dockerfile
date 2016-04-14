FROM alpine:3.3

WORKDIR /go/src/github.com/gaia-adm/mr-burns

COPY .dist/mr-burns /usr/bin/mr-burns

CMD ["--help"]
ENTRYPOINT ["/usr/bin/mr-burns"]