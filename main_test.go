package main

import (
	"fmt"
	"os"
	"testing"
)

// Setup function to create resources before running tests
func setup() {
	fmt.Println("Setting up resources...")

	// Example: Create a directory for testing
	err := os.Mkdir("./.pm/blobs", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating test directory: %v\n", err)
	}

	err = os.Mkdir("./.pm/trie", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating test directory: %v\n", err)
	}

	err = os.Mkdir("./.pm/dag", os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating test directory: %v\n", err)
	}
}

// Teardown function to clean up after tests
func teardown() {
	fmt.Println("Tearing down resources...")

	// Path to your trie directory
	pmDir := ".pm"

	// Remove the entire directory and its contents
	err := os.RemoveAll(pmDir)
	if err != nil {
		fmt.Printf("failed to remove .pm directory: %v", err)
	}
}

// TestMain is the entry point for running tests
func TestMain(m *testing.M) {
	// Call the setup function before running tests
	setup()

	// Run the tests
	exitCode := m.Run()

	// Call the teardown function after running tests
	teardown()

	// Exit with the code returned by m.Run()
	os.Exit(exitCode)
}
