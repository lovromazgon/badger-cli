# badger-cli

[![License](https://img.shields.io/github/license/lovromazgon/badger-cli)](https://github.com/ConduitIO/conduit/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovromazgon/badger-cli)](https://goreportcard.com/report/github.com/lovromazgon/badger-cli)

`badger-cli` is a simple command-line interface for interacting with [Badger DB](https://github.com/dgraph-io/badger), 
a fast key-value database written in Go.

## Features

- Connect to an existing Badger DB
- Get, set, and delete key-value pairs
- List keys with optional glob pattern matching

## Installation

To install `badger-cli`, make sure you have Go installed on your system, then run:

```sh
go install github.com/lovromazgon/badger-cli
```

## Usage

Run the CLI by providing the path to your Badger database:

```sh
badger-cli /path/to/your/badger/db
```

Once the CLI is running, you can use the following commands:

- `get <key>`: Retrieve the value for a given key
- `set <key> <value>`: Set a value for a given key
- `delete <key>`: Delete a key-value pair
- `list [pattern]`: List all keys, optionally filtered by a glob pattern
- `exit`: Exit the CLI

### Examples

```sh
> set mykey myvalue
Value set successfully
> get mykey
myvalue
> list my*
mykey
> delete mykey
Value deleted successfully
> list
No matching keys found
> exit
```

## Glob Pattern Matching

The `list` command supports glob pattern matching:

- `*`: Matches any sequence of characters
- `?`: Matches any single character
- `[abc]`: Matches any character in the set
- `[a-z]`: Matches any character in the range

Examples:
- `list app_*`: Lists all keys starting with "app_"
- `list *_config`: Lists all keys ending with "_config"
- `list user_??`: Lists all keys starting with "user_" followed by exactly two characters

## Acknowledgements

The initial version of this CLI was developed with assistance from an AI language model. The code has since been modified and expanded.

