#!/usr/bin/env bash

export SRCDIR=$GOPATH/src/github.com/nburunova
go version
go env
mkdir -p $SRCDIR
cd $SRCDIR
ln -s /app/currency-loader
cd currency-loader
set -x
if [ $? -ne 0 ] ;
then
    exit 1
fi
make build-api
pwd
