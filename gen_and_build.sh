#!/bin/bash

set -e

go run gen.go
cp index.md doc_source/
mkdocs build
echo "In the case of a failure, ensure that you have run 'workon cas'"
