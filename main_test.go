package main


import (
  "os"
  "fmt"
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

	// Path to your blobs directory
  blobsDir := ".pm/blobs"

  // Remove the entire directory and its contents
  err := os.RemoveAll(blobsDir)
  if err != nil {
    fmt.Printf("failed to remove blobs directory: %v", err)
  }

  // Path to your trie directory
  trieDir := ".pm/trie"

  // Remove the entire directory and its contents
  err = os.RemoveAll(trieDir)
  if err != nil {
    fmt.Printf("failed to remove trie directory: %v", err)
  }

  // Path to your trie directory
  dagDir := ".pm/dag"

  // Remove the entire directory and its contents
  err = os.RemoveAll(dagDir)
  if err != nil {
    fmt.Printf("failed to remove dag directory: %v", err)
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