# ottodocs

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

`ottodocs` is a command-line tool written in Go that uses GPT-3 to automatically generate or add inline documentation for your code. It can parse a git repository or an individual file and create markdown documentation or add inline comments. The tool requires an OpenAI API key to function.

`ottodocs` utilizes the `just` command runner for building and running tasks, making it easy to use and maintain.

## Installation

There are two methods to install `ottodocs` :

1. **Precompiled binaries:** Download the precompiled binaries from the [GitHub releases tab](https://github.com/chand1012/ottodocs/releases).
2. **Build from source:** Clone the repo and build the binary by running the following commands:

```sh
git clone https://github.com/chand1012/ottodocs.git
cd ottodocs
just build # will output binary to bin/otto. Copy the file to a directory in your PATH
```

## Usage

For detailed usage instructions, please refer to the [documentation](docs/otto.md).
