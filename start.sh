#!/bin/sh

git submodule update --init --recursive
rm cryptoassets
go build
nohup ./cryptoassets &
