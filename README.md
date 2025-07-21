# badger-cli

[![License](https://img.shields.io/github/license/lovromazgon/badger-cli)](https://github.com/ConduitIO/conduit/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/lovromazgon/badger-cli)](https://goreportcard.com/report/github.com/lovromazgon/badger-cli)

`badger-cli` is a simple command-line interface for interacting with [Badger DB](https://github.com/dgraph-io/badger), 
a fast key-value database written in Go.

## Features

- Connect to an existing Badger DB
- Get, set, and delete key-value pairs
- List keys with optional glob pattern matching
- Read-only mode for safe database inspection

## Installation

Install using homebrew:

```sh
brew install lovromazgon/tap/badger-cli
```

Or build it from source using Go:

```sh
go install github.com/lovromazgon/badger-cli
```

Or download the binary manually from the [latest release](https://github.com/lovromazgon/badger-cli/releases/latest).

## Usage

Run the CLI by providing the path to your Badger database:

```sh
badger-cli /path/to/your/badger/db
```

### Read-only Mode

To open the database in read-only mode (useful for inspecting production databases safely):

```sh
badger-cli -readonly /path/to/your/badger/db
# or
badger-cli --readonly /path/to/your/badger/db
```

**Note**: The order of arguments is flexible. You can place the readonly flag before or after the database path:

```sh
badger-cli -readonly /path/to/your/badger/db
badger-cli /path/to/your/badger/db -readonly
```

In read-only mode:
- The prompt will show `[READONLY] >`
- Only `get` and `list` commands are available
- `set` and `delete` commands are disabled
- The database is opened with read-only permissions
- The database must already exist (will not create new databases)
- If the database is already in use by another process, you may need to run without -readonly first to ensure proper initialization

Once the CLI is running, you can use the following commands:

- `get <key>`: Retrieve the value for a given key
- `set <key> <value>`: Set a value for a given key (not available in read-only mode)
- `delete <key>`: Delete a key-value pair (not available in read-only mode)
- `list [pattern]`: List all keys, optionally filtered by a glob pattern
- `exit`: Exit the CLI

### Examples

Normal mode:
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

Read-only mode:
```sh
[READONLY] > get mykey
myvalue
[READONLY] > list my*
mykey
[READONLY] > set newkey newvalue
Error: Cannot set values in read-only mode
[READONLY] > delete mykey
Error: Cannot delete values in read-only mode
[READONLY] > exit
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
