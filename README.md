# MQTThing

A command line tool to subscribe to 
[The Things Network](https://www.thethingsnetwork.org) MQTT brokers.

## Installation

A recent version of [Go](https://go.dev) is required to build the program.

Compile the app with `go build .`
and you will get a `mqtthing` binary in this folder.

Run `go install` to install the app in you Go bin folder
or copy the binary to somewhere on your `$PATH`.

## Usage

The app will recognize these flags:

* `-b` Broker to connect to (defaults to `127.0.0.1:1883`)
* `-u` Username
* `-p` Password (API key)
* `-t` Topic (defaults to `#`)
* `-raw` Prints the entire message

Without the `-raw` flag the app will parse the decoded payload and print it.

