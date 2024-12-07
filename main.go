package main

import (
	"fmt"
	"os"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}
func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: my-cli-app <command> [options]")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "login":
		retrieveCreds()
	default:
		fmt.Printf("Unknown command %s\n", os.Args[1])
	}
}

func retrieveCreds() {
	dir, err := os.UserHomeDir()
	check(err)
	fmt.Println(dir)

	dat, err := os.ReadFile(dir + "/.aws/credentials")
	if err != nil {
		fmt.Println("Failed to retrieve data ", err)
	}
	fmt.Printf("Here's the content of the file : %s", dat)
}
