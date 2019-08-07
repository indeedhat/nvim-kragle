#!/bin/bash

CGO_ENABLED=0 go build .
GOOS=darwin go build -o kragle.darwin .
