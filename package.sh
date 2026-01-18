#!/usr/bin/env bash
# Create a distributable zip (age-wrapper.zip) containing the project files.
set -euo pipefail
FILES="main.go util.go encrypt.go decrypt.go keygen.go go.mod README.md LICENSE package.sh"
ZIPNAME=age-wrapper.zip
rm -f $ZIPNAME
echo "Creating $ZIPNAME..."
zip -r $ZIPNAME $FILES
echo "Created $ZIPNAME"