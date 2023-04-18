# OttoDocs 🦦

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

OttoDocs is a command-line tool written in Go that uses GPT-3 (and GPT-4 once the API is available) to automatically generate or add inline and markdown documentation for your code. It can parse a git repository or an individual file and create markdown documentation or add inline comments. The tool requires an [OpenAI API key](https://platform.openai.com/account/api-keys) to function.

OttoDocs utilizes the `just` command runner for building and running tasks, making maintaining the project easier. If you do not have `just` installed, see [here](https://just.systems/man/en/chapter_5.html) for installation methods.

## Installation

There are two methods to install OttoDocs:

1. **Precompiled binaries:** Download the precompiled binaries from the [GitHub releases tab](https://github.com/chand1012/ottodocs/releases).
2. **Build from source:** Clone the repo and build the binary by running the following commands:

```sh
git clone https://github.com/chand1012/ottodocs.git
cd ottodocs
just build # will build to bin/otto
# or
just install # will build & copy to $GOPATH/bin. 
```

If you want to build for all supported platforms, you can run the following command:

```sh
just crossbuild
```

This will build and compress binaries for all supported platforms and place them in the `dist` directory.

## Getting Started

First, you need to create an OpenAI API Key. If you do not already have an OpenAI account, you can create a new one and get some free credits to try it out. Once you have an account, you can create an API key by going to the [API Keys tab](https://platform.openai.com/account/api-keys) in your account settings.

Once you have an API key, you can log in to ottodocs by running the following command:

```sh
otto login
```

Optionally you can pass the API key as an argument to the command:

```sh
otto login --apikey $OPENAI_API_KEY
```

Once that is complete, you can start generating documentation by running the following command:

```sh
otto docs <path to repo or file>
```

Or for a single file, you can run:

```sh
otto doc -f <path to file>
```

## Usage

For detailed usage instructions, please refer to the [documentation](docs/otto.md).
