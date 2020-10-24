#!/bin/bash -e

function clean() {
				find . -iname 'Test*.png' | xargs rm
}

function test() {
				go fmt ./...
				go test ./...
}

COMMAND="${1}"

case "${COMMAND}" in
				clean)
								clean
								;;
				test)
								test
								;;
				*)
								test
								;;
esac
