#!/bin/bash
# Run go test with coverage and output as junit report
workdir=.cover

generate_test_coverage_data() {

    # create dir if not exist and remove old content
    mkdir "$workdir"
	rm -rf "$workdir/*"
    #glide install
    go get
    cd controller
    go test -cover -v > "../$workdir/go-results_tests.txt"
    cat "../$workdir/go-results_tests.txt" | go-junit-report > "../$workdir/go-results_tests.xml"
}

generate_test_coverage_data