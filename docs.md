# justfile
This file contains commands for building, running, and cleaning the project.

## Default
The default command lists all available commands in the file, in an unsorted manner.

## Build
This command builds the project by running `go build`. It also includes the `-v` flag, which prints the names of packages as they are compiled. Upon successful compilation, an executable named `otto` will be created.

## Clean
This command removes the previously built `otto` executable.

## Run
The `run` command is used to execute the `main.go` file with additional arguments. This command takes in variable arguments, i.e., a list of commands that will be passed to `main.go`. The `{{commands}}` placeholder in the command is replaced with the commands passed as arguments when the command is executed.

### Inputs
The commands passed as arguments when executing the command.

### Outputs
The execution of the `main.go` file with the specified commands.# isJustfile

## Package: cmd
This code snippet is part of the `cmd` package.

## Variable: question
`question` is a global variable of type `string`.

## Command: chatCmd
`chatCmd` is a command-line sub-command implemented using `cobra.Command`. It represents a functionality to chat with ChatGPT.

### Use
```
ottodocs chat [flags]
```

### Flags
- `--question`, `-q`: A question to chat ChatGPT.

### Functionality
This command chats with ChatGPT using the OpenAI API. If an API key is not set, it prompts the user to run `ottodocs login` to first login. If no `question` is provided as a flag, the user is asked to input a question. It then sends the question to ChatGPT and displays the response.

### Inputs
- `cmd *cobra.Command` - a pointer to a `cobra.Command` representing the chat sub-command.
- `args []string` - an array of strings representing any flags or extra arguments pass to the chat sub-command.

### Outputs
The output of this command is a string representing the response of ChatGPT to the provided question.# justfile

## Package and Description
The package `cmd` contains the code for handling the command-line interface (CLI) of OttoDocs. `justfile` is a command-line tool for adding documentation to your code using the OpenAI ChatGPT API.

## `docCmd` Variable
`docCmd` is a variable that represents the `doc` command. The `doc` command is used to add documentation to the code. It is a subcommand of the root command.

## `Run` Function
The `Run` function contains the implementation of the `doc` command. It takes in the command and arguments as input and adds documentation to the code using the OpenAI ChatGPT API. 

### Input
- `cmd` is a pointer to `cobra.Command`. It is used to create the `doc` command.
- `args` is a list of strings. It contains the command-line arguments passed to the `doc` command.

### Output
The output of the `Run` function depends on the provided options:
- If the `--output` flag is set and specifies a file path, the function writes the generated documentation to the specified file.
- If the `--overwrite` flag is set, the function overwrites the original file with the generated documentation.
- If neither `--output` nor `--overwrite` flags are set, the function prints the generated documentation to the console.

## `init` Function
The `init` function initializes the `docCmd` command and its flags. The flags are:
- `file`, which specifies the file to add documentation to.
- `prompt`, which specifies the prompt to use when generating the documentation.
- `output`, which specifies the file to write the generated documentation to.
- `inline`, which indicates whether to add the documentation to the code inline.
- `markdown`, which indicates whether to output the documentation in markdown.
- `overwrite`, which indicates whether to overwrite the original file with the generated documentation.# Justfile
## Package cmd
This package cmd provides a Cobra command `docs` which documents an entire repository of files. The `docs` command takes the path to the repository as the first positional argument.

## Var docsCmd
This variable represents the `docs` command.
### Use
`docs`
### Short
Document a repository of files
### Long
Document an entire repository of files. Specify the path to the repo as the first positional argument. This command will recursively search for files in the directory and document them.

## Function init()
This function initializes the `docs` command and its flags.
### Flags
- `prompt`: StringVarP. Specifies the prompt to use for the ChatGPT API.
- `output`: StringVarP. Specifies the path to the output file. For use with the `--markdown` flag.
- `ignore`: StringVarP. Specifies the path to the `.gptignore` file.
- `markdown`: BoolVarP. Indicates the output format is in Markdown.
- `inline`: BoolVarP. Indicates the output format is inline.
- `overwrite`: BoolVarP. Overwrites the original file.
- `ignore-gitignore`: BoolVarP. Ignore the `.gitignore` file.

## Function Run()
This function runs the `docs` command. It takes the following inputs:
- `cmd`: a pointer to a `cobra.Command` object.
- `args[] string`: slice of repository paths to be processed.

### Functionality
The function checks for several errors involved with the execution of the `docs` command. It then process the given repository and then calls `document.SingleFile` function, passing the file path, chatPrompt and API Key from retrieved configuration. The `document.SingleFile` function returns a string containing the processed and documented file which is written to either the output file, if specified, or printed to stdout.

### Inputs
- `cmd *cobra.Command`: Command Object
- `args []string`: Slice of input arguments

### Outputs
The `Run()` function produces a string that contains the documentation of given repository files.# Justfile

## Package cmd

This package provides a CLI app to interact with the OpenAI ChatGPT API. This specific file, named "Justfile", defines the "login" subcommand.

## loginCmd

### Use

```
login
```
  
### Short

Add an API key to your configuration

### Long

Add an API key to your configuration. This API key will be used to authenticate with the OpenAI ChatGPT API.

### Run

This function is executed when the "login" command is called. It first checks whether an `apiKey` variable was passed as a flag or not. If not, it prompts the user to provide it through the terminal. Then, it saves the `apiKey` to a configuration file located at `~/.ottodocs/config.json` using the `config.Load()` and `config.Save()` functions. If an error occurs during the configuration file's manipulation, an error message is printed to the terminal, and the application exits with an error code. If all goes well, a message is printed to the terminal confirming the successful save of the `apiKey`.

## init()

This function adds the "login" subcommand to the root command created in another file, "rootCmd". It also defines the "apikey" flag using the `loginCmd.Flags()` function.# isJustfile

## `cmd` package

This package contains an implementation of a command-line interface (CLI) that generates a ChatGPT prompt from a given Git repository.

The package is imported as:

```go
import (
  "github.com/chand1012/git2gpt/prompt"
  "github.com/spf13/cobra"
)
```

### `promptCmd`

`promptCmd` is a global variable of type `*cobra.Command` that represents the main prompt command. It has the following properties:

- `Use`: a string representing the command name, "prompt".
- `Short`: a string containing a brief one-line description of the command.
- `Long`: a string containing a more detailed description of the command.
- `Args`: a `cobra.PositionalArgs` function that validates that the user provides a path to a Git repository as the first positional argument.
- `Run`: a function that runs the prompt command.

### `init()`

`init()` is a function that initializes the `promptCmd` variable and adds it to the `rootCmd`, which is a top-level command that contains all other commands in the CLI.

### `Flags`

The `promptCmd` command has several flags that can be used to customize the behavior of the command. These are:

- `preambleFile`: a flag of type `string` that specifies the path to a preamble text file. The flag name is `-p` or `--preamble`.
- `outputFile`: a flag of type `string` that specifies the path to the output file. The flag name is `-o` or `--output`.
- `estimateTokens`: a flag of type `bool` that indicates whether to estimate the number of tokens in the output. The flag name is `-e` or `--estimate`.
- `ignoreFilePath`: a flag of type `string` that specifies the path to a `.gptignore` file. The flag name is `-i` or `--ignore`.
- `ignoreGitignore`: a flag of type `bool` that indicates whether to ignore the `.gitignore` file. The flag name is `-g` or `--ignore-gitignore`.
- `outputJSON`: a flag of type `bool` that indicates whether to output in JSON format. The flag name is `-j` or `--json`.# Justfile

## Package cmd
This package provides a CLI interface to run commands relating to documentation. 

## rootCmd 
This variable represents the base command when it is called without any subcommands. 

### Inputs
- Use: (string) The name of the command used in the CLI interface.
- Short: (string) A brief description of what the command does.
- Long: (string) A more detailed description of what the command does. 

## Execute
This function is called by main.main() and adds all child commands to the root command while setting flags appropriately. 

### Outputs
- Error: (error) Any errors encountered while executing the command. 

## init 
This function defines flags and configuration settings that will be global for the application. 

### Persistent Flags
- config (string): This specifies the path to the configuration file. 

### Local Flags
- toggle (bool): This is a local flag that's only run when the action is called directly. This flag can be helpful during debugging.# isJustfile

## Variables

### repoPath
- Type: string
- Description: holds the path of the repository

### preambleFile
- Type: string
- Description: holds the path of the preamble file

### outputFile
- Type: string
- Description: holds the path of the output file

### estimateTokens
- Type: bool
- Description: If set to true, estimates the token count while processing the file.

### ignoreFilePath
- Type: string
- Description: holds the path of the file containing patterns to ignore

### ignoreGitignore
- Type: bool
- Description: If set to true, ignores files specified by gitignore.


### outputJSON
- Type: bool
- Description: If set to true, the output file will be in JSON format.


### filePath
- Type: string
- Description: holds the path of the file


### chatPrompt
- Type: string
- Description: holds the prompt message for chatbot


### inlineMode
- Type: bool
- Description: If set to true, displays the output in inline mode. 


### markdownMode
- Type: bool
- Description: If set to true, displays the output in markdown format.


### overwriteOriginal
- Type: bool
- Description: If set to true, the original file will be overwritten with the output.


## Functions and Classes

None. This file only contains variable declarations that are used across multiple files in the package.# config

## Config 

`Config` represents a struct that holds the configuration file. It has only one field, `APIKey`, with type `string`.

## createIfNotExists

`createIfNotExists` is a function that returns the path to the `config.json` file and an error. It checks if the `~/.ottodocs/config.json` file exists. If it does not exist, it creates a new one, and adds an empty JSON object `{}` to the file. It returns the path to the config file and any errors encountered.

## Load

`Load` is a function that loads the configuration from the `config.json` file located in the folder `~/.ottodocs/`. If the file or folder does not exist, `Load` creates them. It returns a pointer to the loaded `Config` object and any errors encountered.

## Save

`Save` is a method of the `Config` type. It takes no arguments and saves the current `Config` object to the `config.json` file located in the folder `~/.ottodocs/`. If the file or folder does not exist, `Save` creates them. It returns any errors encountered.# constants

## CommentOperators

This variable is a map that contains the comment highlights (comment operators) of various programming languages. It maps the file extension of a programming language to its comment operator.

### Input

This variable has no input.

### Output

This variable returns a map containing the comment operator for each supported programming language. The key of each pair is the file extension of the language, and the corresponding value is the operator that identifies a comment in that language. The following programming languages are currently supported: Python, Go, C, C++, C/C++ Header, C#, Java, JavaScript, TypeScript, PHP, Ruby, Rust, Swift, Shell Script, Perl, Lua, MATLAB, R, Scala, Kotlin, Visual Basic .NET, Fortran, Assembly, HTML and CSS.# isJustfile
## constants
### OPENAI_MAX_TOKENS
This constant is an integer and has a value of 4000. It is defined in the constants package. It can be used to set the maximum number of tokens to be used in OpenAI's GPT models.# constants

This package provides two string constants that contain prompts for assisting with code documentation.

## DOCUMENT_FILE_PROMPT

This string constant provides a prompt for OpenAI API when documenting a file. It contains instructions for how to document code, including the rules for line numbers, avoiding the use of ranges, and keeping documentation in English. It also emphasizes that each line must be documented with a single string and in the order of the code.

**Input:** None

**Output:** This constant is a string.

## DOCUMENT_MARKDOWN_PROMPT

This string constant provides a prompt for documenting code in valid markdown. It contains instructions for how to document the name of a file, functions, and classes. It also emphasizes the importance of describing the use, inputs, and outputs of the code and each of its public functions. Like DOCUMENT_FILE_PROMPT, it reiterates that documentation must be in the order of the code and must not contain author or copyright information.

**Input:** None

**Output:** This constant is also a string.# isJustfile

## Markdown
This function takes in a `filePath` string, `chatPrompt` string, and `APIKey` string as inputs, and returns a `string` and an `error` as outputs. The `filePath` specifies the path to a markdown file that is going to be parsed, and `chatPrompt` is a string which will be used to prompt the OpenAI Chatbot, which will be used to parse the markdown file. `APIKey` is an authentication key required to use the OpenAI API.

```go
func Markdown(filePath string, chatPrompt string, APIKey string) (string, error)
```

### Inputs
- `filePath` (`string`) : The location of the markdown file to be parsed.
- `chatPrompt` (`string`) : The prompt to be used to converse with the OpenAI chatbot.
- `APIKey` (`string`) : The OpenAI API authentication key.

### Outputs
- `string` : The processed markdown file that the OpenAI chatbot generated as parsed output.
- `error` : An error message, in case an error occurs while parsing the file.

## Other Functions and Packages

### Function: None

### Package: `document`

This package imports various Go and third-party packages to parse markdown files and use OpenAI's chatbot, which helps to generate parsed markdown files. There is no use of this package outside this file. The imported packages are:

- `os` : for file read operations.
- `github.com/CasualCodersProjects/gopenai` : an open-source package to interact with OpenAI's API.
- `github.com/CasualCodersProjects/gopenai/types` : this package has OpenAI request and response types.
- `github.com/chand1012/ottodocs/constants` : this package has some constants used within the application.
- `github.com/pandodao/tokenizer-go` : this package is used to calculate the number of tokens from the chatbot message.# isJustfile

## extractLineNumber

This function takes a string input `line` and returns an `int` and an `error`. If `line` does not contain an "-" character, `extractLineNumber` parses the string into an integer and returns it. Otherwise, `extractLineNumber` returns the first integer obtained by splitting `line` at the "-" character.

### Inputs
- line: `string` 

### Outputs
- int: `int` - Line number obtained from parsing `line`
- error: `error` - Error message, if any.

## SingleFile

This function takes three string inputs:
- `filePath` - file path of the file to be documented
- `chatPrompt` - prompt for the AI to generate documentation for the file
- `APIKey` - an OpenAI API key

`SingleFile` documents a file using the OpenAI ChatGPT API. The function returns a string that represents the file with further documentation inserted. If an error occurs while processing, an error message is returned describing what went wrong.

### Inputs
- filePath: `string` - path of file to be documented
- chatPrompt: `string` - prompt for the AI to generate documentation for the file
- APIKey: `string` - an OpenAI API key 

### Outputs
- string: `string` - Contents of the file with further documentation inserted.
- error: `error` - Error message, if any.# justfile

## Description 
This file is used to declare the dependencies needed for the project.

## Dependencies
The following dependencies are required for the project:
- github.com/CasualCodersProjects/gopenai v0.3.0
- github.com/chand1012/git2gpt v0.3.0
- github.com/pandodao/tokenizer-go v0.1.0
- github.com/spf13/cobra v1.6.1

Additionally, the following indirect dependencies are also required:
- github.com/dlclark/regexp2 v1.8.1 
- github.com/dop251/goja v0.0.0-20230304130813-e2f543bf4b4c 
- github.com/dop251/goja_nodejs v0.0.0-20230226152057-060fa99b809f 
- github.com/go-sourcemap/sourcemap v2.1.3+incompatible 
- github.com/gobwas/glob v0.2.3 
- github.com/google/pprof v0.0.0-20230309165930-d61513b1440d 
- github.com/inconshreveable/mousetrap v1.0.1 
- github.com/spf13/pflag v1.0.5 
- golang.org/x/text v0.8.0 

Note that these dependencies are downloaded automatically by the Go package manager when the program is built or run.# isJustfile

## main function
Entry point of the program. Calls the `cmd.Execute()` function.

### Inputs
No inputs.

### Outputs
No outputs.

## Package
- `github.com/chand1012/ottodocs/cmd`

This program imports the `cmd` package from the `github.com/chand1012/ottodocs` repository.

## Dependencies
None

This program has no external dependencies other than the `cmd` package which is imported from a different repository.# textfile package
This package contains functions for text operations.

## InsertLine
This function inserts a new line into a string of text on the specified line number.

### Inputs
* `code` - a string of text to be inserted into.
* `lineNumber` - an integer representing the line number to be inserted at.
* `newText` - the text to be inserted.

### Outputs
* `string` - the updated text string with the new line inserted.
* `error` - an error message if line number is less than 1 or greater than the number of lines in the text string.

## InsertLinesAtIndices
This function inserts multiple new lines into a string of text at specified indices.

### Inputs
* `file` - a string of text to be inserted into.
* `indices` - a slice of integers representing the line numbers to be inserted at.
* `linesToInsert` - a slice of strings containing the text to be inserted at each index.

### Outputs
* `string` - the updated text string with the new lines inserted.
* `error` - an error message if the length of indices and linesToInsert isn't equal, or if any index is less than 1 or greater than the number of lines in the text string.