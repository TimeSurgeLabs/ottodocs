default:
  just --list --unsorted

build:
  go build -v -o bin/otto

clean:
  rm -f otto

run *commands:
  go run main.go {{commands}}

cobra-docs:
  go run docs/gen_docs.go

install: build
  cp bin/otto $GOPATH/bin

crossbuild:
  #!/bin/bash

  # Set the name of the output binary and Go package
  BINARY_NAME="otto"
  GO_PACKAGE="github.com/chand1012/ottodocs"

  # Build for M1 Mac (Apple Silicon)
  echo "Building for M1 Mac (Apple Silicon)"
  env GOOS=darwin GOARCH=arm64 go build -o "${BINARY_NAME}" "${GO_PACKAGE}"
  zip "${BINARY_NAME}_darwin_arm64.zip" "${BINARY_NAME}"
  rm "${BINARY_NAME}"

  # Build for AMD64 Mac (Intel)
  echo "Building for AMD64 Mac (Intel)"
  env GOOS=darwin GOARCH=amd64 go build -o "${BINARY_NAME}" "${GO_PACKAGE}"
  zip "${BINARY_NAME}_darwin_amd64.zip" "${BINARY_NAME}"
  rm "${BINARY_NAME}"

  # Build for AMD64 Windows
  echo "Building for AMD64 Windows"
  env GOOS=windows GOARCH=amd64 go build -o "${BINARY_NAME}.exe" "${GO_PACKAGE}"
  zip "${BINARY_NAME}_windows_amd64.zip" "${BINARY_NAME}.exe"
  rm "${BINARY_NAME}.exe"

  # Build for AMD64 Linux
  echo "Building for AMD64 Linux"
  env GOOS=linux GOARCH=amd64 go build -o "${BINARY_NAME}" "${GO_PACKAGE}"
  tar czvf "${BINARY_NAME}_linux_amd64.tar.gz" "${BINARY_NAME}"
  rm "${BINARY_NAME}"

  # Build for ARM64 Linux
  echo "Building for ARM64 Linux"
  env GOOS=linux GOARCH=arm64 go build -o "${BINARY_NAME}" "${GO_PACKAGE}"
  tar czvf "${BINARY_NAME}_linux_arm64.tar.gz" "${BINARY_NAME}"
  rm "${BINARY_NAME}"
