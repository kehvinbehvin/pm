package main

import (
	"fmt"
	"os"

	"encoding/gob"
	"golang.org/x/crypto/ssh"

	"github/pm/dag"
	"github/pm/resolver"
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

func pushDeltas(localTree *dag.DeltaTree, remoteTree *dag.DeltaTree) {
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

	// Local should be ahead here.
	// If remote is ahead or is there is a deviation, we should abort
	lastCommonHash, longerTree, _, lcsType, lcsErr := resolver.CalculateLcs(localTree, remoteTree)
	if lcsErr != nil {
		fmt.Println("LCS Err")
		return
	}

	if lcsType != dag.LocalAhead {
		fmt.Println("LCS not localahead")
		return
	}
	deltasToPush, deltaErr := resolver.GetDeltasAhead(longerTree, lastCommonHash)
	if deltaErr != nil {
		fmt.Println("Deltas ahead error")
		return
	}

	cmd := "put delta" // Example command to save data to a file
	if err := session.Start(cmd); err != nil {
		fmt.Printf("Failed to start command: %v", err)
	}

	if srcStat.Size() > 0 {
		gob.Register(&dag.VertexDelta{})
		gob.Register(&dag.EdgeDelta{})

		encoder := gob.NewEncoder(w)
		encodingErr := encoder.Encode(deltasToPush)
		if encodingErr != nil {
			fmt.Printf(encodingErr.Error())
			fmt.Println("Error encoding delta")
			return
		}
		fmt.Println("Deltas pushed")
	}
}
