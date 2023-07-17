# OttoDocs 🦦

[![License: AGPL v3](https://img.shields.io/badge/License-AGPL%20v3-blue.svg)](https://www.gnu.org/licenses/agpl-3.0)

Ottodocs is a command-line tool built in Go that utilizes OpenAI's ChatGPT API to generate commit messages, pull requests, answers to questions, and even shell commands based on input prompts and code context. It helps developers automate various parts of their development workflow using AI. The tool requires an [OpenAI API key](https://platform.openai.com/account/api-keys) to function.

## Installation

There are two methods to install OttoDocs:

1. **Precompiled binaries:** Download the precompiled binaries from the [GitHub releases tab](https://github.com/TimeSurgeLabs/ottodocs/releases).
2. **Build from source:** Clone the repo and build the binary by running the following commands:

### Installing Precompiled Binaries

Simply download the binary for your platform from the [GitHub releases tab](https://github.com/TimeSurgeLabs/ottodocs/releases) and place it in your `$PATH` .

### Building From Source

OttoDocs utilizes the `just` command runner for building and running tasks, making maintaining the project easier. If you do not have `just` installed, see [here](https://just.systems/man/en/chapter_5.html) for installation methods. Ottodocs requires [Go 1.20+](https://go.dev/dl/) to build.

```sh
git clone https://github.com/TimeSurgeLabs/ottodocs.git
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
otto config --apikey $OPENAI_API_KEY
```

You can set the model to use by running:

```sh
otto config --model $MODEL_NAME
```

You can add a GitHub Personal Access Token for opening PRs by running:

```sh
otto config --token $GITHUB_TOKEN # Optional. Only needed for opening PRs and reading Issues
```

Make sure that your access token has the `repo` scope.

Once that is complete, you can start running commands!

## Usage

More detailed usage can be found in the [documentation](https://ottodocs.chand1012.dev/docs/usage/otto).

### Chat

![Made with VHS](https://vhs.charm.sh/vhs-7lC4zu09dmW4TZFOyQtCKI.gif)

To Chat with ChatGPT from the commandline, use chat.

```
otto chat
```

### Code Generation

![Made with VHS](https://vhs.charm.sh/vhs-5A6u3ITYSIp2qd1T7XdYEM.gif)

Otto can generate code and save it directly to your project! To do this, you can run the following command:

```sh
otto edit <path to file> -g "Write me a Python function that returns the sum of two numbers"
```

You can also give Otto additional context to help it generate better code:

```sh
otto edit <path to file> -g "Write me a Python function that returns the sum of two numbers" -c hello_world.py
```

Or you can have it use as much of the repo as possible as context:

```sh
otto edit <path to file> -g "Write me a Python function that returns the sum of two numbers" -r
```

### Documentation

```sh
otto docs <path to repo or file>
```

Or for a single file, you can run:

```sh
otto doc -f <path to file>
```

### Ask

Ask a question about a repo:

```sh
otto ask . -q "What does LoadFile do differently than ReadFile?"
```

### Commit Messages

Generate a commit message:

![Made with VHS](https://vhs.charm.sh/vhs-4Uti5pLyUQ85pueoJH5IQ.gif)

```sh
otto commit # optionally add --push to push to remote
```

### Pull Request

Generate a pull request:

```sh
# make sure you are creating the PR on the correct base branch
otto pr -b main # optionally add --publish to publish the Pull Request
```

### Release Notes

Generate release notes:

```sh
otto release # optionally add --publish to publish the release
```

### Command Ask

Ask it about commands:

```sh
otto cmd -q "what is the command to add a remote?"
```

## Usage

For detailed usage instructions, please refer to the [documentation](https://ottodocs.chand1012.dev/docs/usage/otto).
