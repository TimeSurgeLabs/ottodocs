default:
  just --list --unsorted

build:
  go build -v -o otto

clean:
  rm -f otto

run *commands:
  go run main.go {{commands}}

cobra-docs:
  go run docs/gen_docs.go
