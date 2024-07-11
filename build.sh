#!/bin/bash

rm -fr ./build/
rm -f ./user_faker.zip
cp -r ./build_template/ ./build/
GIN_MODE=release GOOS=darwin GOARCH=arm64 go build -o ./build/run_arm64
GIN_MODE=release GOOS=darwin GOARCH=amd64 go build -o ./build/run_amd64

lipo -create -output ./build/run ./build/run_arm64 ./build/run_amd64
rm -f ./build/run_*

zip -r user_faker.zip ./build