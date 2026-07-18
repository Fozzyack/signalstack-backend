package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password, err := passwordInput()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to generate password hash:", err)
		os.Exit(1)
	}

	fmt.Println(string(hash))
}

func passwordInput() (string, error) {
	if len(os.Args) > 1 {
		return os.Args[1], nil
	}

	fmt.Fprint(os.Stderr, "Password: ")
	scanner := bufio.NewScanner(os.Stdin)
	if !scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		return "", fmt.Errorf("password is required")
	}

	return scanner.Text(), nil
}
