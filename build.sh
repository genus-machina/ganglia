#!/bin/bash -e

function clean() {
	find . -iname 'Test*.png' | xargs rm
}

function test() {
	go fmt ./...
	go test ./...
}

function update() {
	go get -u ./...
	go mod tidy
}

COMMAND="${1}"

case "${COMMAND}" in
	clean)
		clean
		;;
	test)
		test
		;;
	update)
		update
		;;
	*)
		test
		;;
esac
