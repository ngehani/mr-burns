FROM golang:1.6-alpine

WORKDIR /go/src/github.com/gaia-adm/mr-burns

# install Git apk
RUN apk --update add git bash \
    && rm -rf /var/lib/apt/lists/* \
    && rm /var/cache/apk/*

# install glide package manager
RUN curl -Ls https://github.com/Masterminds/glide/releases/download/0.10.1/glide-0.10.1-linux-amd64.tar.gz | tar xz -C /tmp \
 && mv /tmp/linux-amd64/glide /usr/bin/

# gox - Go cross compile tool
# cover - Go code coverage tool
# go-junit-report - convert Go test into junit.xml format
RUN go get github.com/mitchellh/gox \
    && go get github.com/jstemmer/go-junit-report

ENV RESULT_DIR /src/.cover
ENV RESULT_FILE go-results_tests.xml

LABEL test=
LABEL test.results.dir=$RESULT_DIR
LABEL test.results.file=$RESULT_FILE
LABEL test.cmd=script/go_test.sh

RUN chmod u+x script/go_build.sh
RUN chmod u+x script/go_test.sh

CMD ["script/go_build.sh"]
