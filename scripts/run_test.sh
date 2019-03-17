#!/bin/sh

set -e

git remote set-url origin https://github.com/mas9612/nwspeaker.git

go test -v ./...
