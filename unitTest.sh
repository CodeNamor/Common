#!/usr/bin/env bash

PRE=$(go test -covermode=count -coverprofile=count.txt ./...)
COVERAGE=$(go tool cover -func=./count.txt)

#The statements below are to extract the percentage of code coverage so that we can use if logic later on to assert X amount of coverage
#PERCENT=${COVER#*\)}
#PERCENT=${PERCENT%%.*}

echo ${COVERAGE}