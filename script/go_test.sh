#!/bin/sh
# Run go test with coverage and output as junit report
workdir=.cover

generate_test_coverage_data() {

    # create dir if not exist and remove old content
    mkdir "$workdir"
	rm -rf "$workdir/*"
    cd dockerclient
    go get
    go test -cover -v | go-junit-report > "../$workdir/go-results_tests.xml"
}

generate_test_coverage_data