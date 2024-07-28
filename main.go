package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/gobwas/glob"
)

func main() {
	// Set default path and check for command-line argument
	var dbPath string
	if len(os.Args) > 1 {
		dbPath = os.Args[1]
	}
	if dbPath == "" {
		log.Fatal("Please supply a path to a badger database")
	}

	// Open the Badger database
	db, err := badger.Open(badger.DefaultOptions(dbPath))
	if err != nil {
		log.Fatalf("Error opening database: %v\n", err)
	}
	defer db.Close()

	// Start the CLI loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" {
			break
		}

		handleCommand(db, input)
	}
}

func handleCommand(db *badger.DB, input string) {
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
		if len(args) != 2 {
			fmt.Println("Usage: set <key> <value>")
			return
		}
		setValue(db, args[0], args[1])
	case "delete":
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
		fmt.Println("Unknown command. Available commands: get, set, delete, list")
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
