package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/dgraph-io/badger/v4"
	"github.com/gobwas/glob"
)

func main() {
	// Parse command line arguments
	var readonly bool
	var dbPath string

	args := os.Args[1:] // Skip program name

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-readonly" || arg == "--readonly" {
			readonly = true
		} else if !strings.HasPrefix(arg, "-") {
			// First non-flag argument is the database path
			if dbPath == "" {
				dbPath = arg
			}
		}
	}

	if dbPath == "" {
		log.Fatal("Please supply a path to a badger database")
	}

	// Open the Badger database
	opts := badger.DefaultOptions(dbPath)
	if readonly {
		opts.ReadOnly = true
	}

	db, err := badger.Open(opts)
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Create a custom completer
	completer := createDatabaseCompleter(db, readonly)

	// Set up readline with the custom completer
	prompt := "> "
	if readonly {
		prompt = "[READONLY] > "
	}

	rl, err := readline.NewEx(&readline.Config{
		Prompt:       prompt,
		AutoComplete: completer,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer rl.Close()

	// Start the CLI loop
	for {
		input, err := rl.Readline()
		if err != nil { // io.EOF, readline.ErrInterrupt
			break
		}
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		handleCommand(db, input, readonly)
	}
}

// createDatabaseCompleter creates a dynamic completer for database keys
func createDatabaseCompleter(db *badger.DB, readonly bool) readline.AutoCompleter {
	items := []readline.PrefixCompleterInterface{
		readline.PcItem("get",
			readline.PcItemDynamic(func(line string) []string {
				return getDatabaseKeys(db, line)
			}),
		),
		readline.PcItem("list"),
	}

	// Only add write commands if not in readonly mode
	if !readonly {
		items = append(items, readline.PcItem("set",
			readline.PcItemDynamic(func(line string) []string {
				return getDatabaseKeys(db, line)
			}),
		))
		items = append(items, readline.PcItem("delete",
			readline.PcItemDynamic(func(line string) []string {
				return getDatabaseKeys(db, line)
			}),
		))
	}

	return readline.NewPrefixCompleter(items...)
}

// getDatabaseKeys retrieves keys from the database for autocomplete
func getDatabaseKeys(db *badger.DB, line string) []string {
	// Store found keys
	var keys []string

	// Retrieve keys from the database
	err := db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		// Convert current line to a prefix for matching
		prefix := []byte(line[strings.LastIndexAny(line, " ")+1:])

		for it.Rewind(); it.Valid(); it.Next() {
			key := it.Item().Key()

			// Only add keys that start with the current prefix
			if bytes.HasPrefix(key, prefix) {
				keys = append(keys, string(key))
			}
		}
		return nil
	})

	if err != nil {
		return nil
	}

	return keys
}

func handleCommand(db *badger.DB, input string, readonly bool) {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return
	}

	command := parts[0]
	args := parts[1:]

	switch command {
	case "get":
		if len(args) != 1 {
			fmt.Println("Usage: get <key>")
			return
		}
		getValue(db, args[0])
	case "set":
		if readonly {
			fmt.Println("Error: Cannot set values in read-only mode")
			return
		}
		if len(args) != 2 {
			fmt.Println("Usage: set <key> <value>")
			return
		}
		setValue(db, args[0], args[1])
	case "delete":
		if readonly {
			fmt.Println("Error: Cannot delete values in read-only mode")
			return
		}
		if len(args) != 1 {
			fmt.Println("Usage: delete <key>")
			return
		}
		deleteValue(db, args[0])
	case "list":
		pattern := "*"
		if len(args) > 0 {
			pattern = args[0]
		}
		listKeys(db, pattern)
	default:
		if readonly {
			fmt.Println("Unknown command. Available commands: get, list")
		} else {
			fmt.Println("Unknown command. Available commands: get, set, delete, list")
		}
	}
}

func listKeys(db *badger.DB, pattern string) {
	g, err := glob.Compile(pattern)
	if err != nil {
		fmt.Printf("Invalid pattern: %v\n", err)
		return
	}

	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		found := false
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if g.Match(string(k)) {
				found = true
				fmt.Println(string(k))
			}
		}
		if !found {
			fmt.Println("No matching keys found")
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Error listing keys: %v\n", err)
	}
}

func getValue(db *badger.DB, key string) {
	err := db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			fmt.Println(string(val))
			return nil
		})
	})
	if err != nil {
		fmt.Printf("Error getting value: %v\n", err)
	}
}

func setValue(db *badger.DB, key, value string) {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		fmt.Printf("Error setting value: %v\n", err)
	} else {
		fmt.Println("Value set successfully")
	}
}

func deleteValue(db *badger.DB, key string) {
	err := db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
	if err != nil {
		fmt.Printf("Error deleting value: %v\n", err)
	} else {
		fmt.Println("Value deleted successfully")
	}
}
