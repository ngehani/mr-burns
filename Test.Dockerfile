FROM golang:1.6-alpine

# install Git apk
RUN apk --update add git bash \
    && rm -rf /var/lib/apt/lists/* \
    && rm /var/cache/apk/*

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

CMD ["script/go_test.sh"]
