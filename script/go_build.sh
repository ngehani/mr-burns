#!/bin/bash
distdir=.dist

go_build() {
  
  rm -rf "${distdir}"
  mkdir "${distdir}"
  cd main
  go get
  go build -v -o ${distdir}/mr-burns/main
}

go_build