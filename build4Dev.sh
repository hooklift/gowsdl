#!/usr/bin/env bash
export GHACCOUNT=hooklift
export NAME=gowsdl
export VERSION=v0.2.1

go build -o build/gowsdl -ldflags="-s -w" cmd/gowsdl/main.go

