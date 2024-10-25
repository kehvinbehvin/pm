package main

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/ssh"
)

func retrieveFile() {
	// Load private key
	keyPath := "/Users/kevin/.ssh/local/pm-server/local_rsa"
	privateKey, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// Create the signer for the private key
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only. Replace with secure method.
	}

	// Connect to SSH server
	conn, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		fmt.Println("Error connecting")
	}
	defer conn.Close()

	// Create a session for running the command
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Error creating new session")
	}
	defer session.Close()

	// Execute the command
	cmd := "get delta"
	output, err := session.CombinedOutput(cmd)
	if err != nil {
		fmt.Printf("Failed to run command: %v", err)
	}

	file, err := os.Create("./.pm/remote/delta")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	err = os.WriteFile("./.pm/remote/delta", []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing to file")
		return
	}
}

func putFile() {
	// Load private key
	keyPath := "/Users/kevin/.ssh/local/pm-server/local_rsa"
	privateKey, err := os.ReadFile(keyPath)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// Create the signer for the private key
	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	// Create SSH client configuration
	config := &ssh.ClientConfig{
		User: "user",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only. Replace with secure method.
	}

	// Connect to SSH server
	conn, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		fmt.Println("Error connecting")
	}
	defer conn.Close()

	// Create a session for running the command
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Error creating new session")
	}
	defer session.Close()

	w, err := session.StdinPipe()
	if err != nil {
		return
	}

	src, err := os.Open("./.pm/delta")
	if err != nil {
		return
	}
	defer src.Close()

	srcStat, err := os.Stat("./.pm/delta")
	if err != nil {
		panic(err)
	}

	cmd := "put delta" // Example command to save data to a file
	if err := session.Start(cmd); err != nil {
		fmt.Printf("Failed to start command: %v", err)
	}

	if srcStat.Size() > 0 {
		n, err := io.Copy(w, src)
		if err != nil {
			return
		}

		fmt.Println(n)
	}
}
