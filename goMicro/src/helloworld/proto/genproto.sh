#!/bin/bash
protoc --micro_out=. --go_out=. *.proto
