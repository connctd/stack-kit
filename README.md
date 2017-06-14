# stack-kit [![Build Status](https://travis-ci.org/connctd/stack-kit.svg?branch=master)](https://travis-ci.org/connctd/stack-kit) [![GoDoc](https://godoc.org/github.com/go-kit/kit/log?status.svg)](https://godoc.org/github.com/connctd/stack-kit) [![Coverage Status](https://coveralls.io/repos/github/connctd/stack-kit/badge.svg)](https://coveralls.io/github/connctd/stack-kit)

This repository provides helper libraries to use go-kit in combination with 
[stackdriver](https://cloud.google.com/stackdriver/).

## logging

The [logging](https://godoc.org/github.com/connctd/stack-kit/logging) package provides utilities for go-kit 
logger so stackriver can correctly parse the received logs.
It is always assumed that the JSONLogger will be used. The logging package currently helps with:

* setting the correct severity
* correctly formatted error reports
