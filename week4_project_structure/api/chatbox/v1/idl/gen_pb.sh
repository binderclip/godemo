#!/bin/sh
protoc -I=./ --go_out=../ --go_opt=paths=source_relative ./message.proto