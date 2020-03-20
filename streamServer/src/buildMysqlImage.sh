#!/bin/bash

VERSION=v`cat Version`

docker build --no-cache -t aaronrootanderson:${VERSION}-`git rev-parse --short HEAD` .