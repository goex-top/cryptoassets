#!/bin/sh

echo 'update submodule'
git submodule update --init --recursive
echo 'remove older binary file'
rm cryptoassets
echo 'rebuild project'
go build
echo 'run it with nohup'
nohup ./cryptoassets &
