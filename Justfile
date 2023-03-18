default:
  just --list --unsorted

build:
  go build -v -o otto

clean:
  rm -f otto

run *commands:
  go run main.go {{commands}}
