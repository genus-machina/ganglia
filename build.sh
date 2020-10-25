#!/bin/bash -e

function clean() {
	find . -iname 'Test*.png' | xargs rm
}

function fail() {
	exit 1
}

function test() {
	go fmt ./...
	go test ./...
}

function unknown() {
	echo "Unknown command '${COMMAND}'." >&2
}

function update() {
	GOPROXY=direct go get -u ./...
	go mod tidy
}

function usage() {
	echo "${0}: [command]"
	echo
	echo "Commands:"
	echo -e "\tclean"
	echo -e "\ttest"
	echo -e "\tupdate"
	echo
}

COMMAND="${1}"

case "${COMMAND}" in
	"")
		test
		;;
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
		unknown
		echo
		usage
		fail
		;;
esac
