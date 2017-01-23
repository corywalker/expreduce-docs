#!/bin/bash

set -e

go run gen.go
cp index.md doc_source/
mkdocs build
