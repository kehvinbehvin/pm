package main

import (
	"fmt"
	"os"

	"encoding/gob"
	"golang.org/x/crypto/ssh"

	"github/pm/dag"
	"github/pm/resolver"
)

func retrieveFile(localTree *dag.DeltaTree, remoteTree *dag.DeltaTree) {
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
		User: "kevin@example.com",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only. Replace with secure method.
	}

	// Connect to SSH server
	conn, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		fmt.Println("Error connecting")
		return
	}
	defer conn.Close()

	// Create a session for running the command
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Error creating new session")
		return
	}
	defer session.Close()

	cmd := "get"

	if len(remoteTree.Seq) != 0 {
		lastRemoteDeltaPtr := remoteTree.Seq[remoteTree.Pointer];
		if lastRemoteDeltaPtr == nil {
			fmt.Println("Cannot find last delta on remote")
			return 
		}
		lastRemoteDelta := *lastRemoteDeltaPtr
		lastRemoteHash := lastRemoteDelta.GetId()

		// Execute the command
		cmd = cmd + " " + lastRemoteHash 
	}
	
	stdOut, pipeErr := session.StdoutPipe();
	if pipeErr != nil {
		fmt.Println("Pipe err")
		return
	}

	err = session.Run(cmd)
	if err != nil {
		fmt.Printf("Failed to run command: %v", err)
	}

	// Initialize decoding from the SSH session input (s)
	decoder := gob.NewDecoder(stdOut)

	// Register types to decode correctly
	gob.Register(&dag.VertexDelta{})
	gob.Register(&dag.EdgeDelta{})

	// Decode the incoming data into a slice of Deltas (or whatever structure is expected)
	var deltasToApply []dag.Delta
	if err := decoder.Decode(&deltasToApply); err != nil {
		fmt.Println("Error Decoding", err.Error())
		fmt.Println("Could not decode remote deltas")
		return
	}

	for _, delta := range deltasToApply {
		delta.SetDeltaTree(remoteTree);
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
		User: "kevin@example.com",
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // For testing only. Replace with secure method.
	}

	// Connect to SSH server
	conn, err := ssh.Dial("tcp", "localhost:2222", config)
	if err != nil {
		fmt.Println("Error connecting")
		return
	}
	defer conn.Close()

	// Create a session for running the command
	session, err := conn.NewSession()
	if err != nil {
		fmt.Println("Error creating new session")
		return
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
	var deltasToPush []*dag.Delta

	// Push the common delta as first delta for server to check for conflicts
	commonDelta := longerTree.IdTree[lastCommonHash]
	deltasToPush = append(deltasToPush, commonDelta)

	deltasAhead, AheadErr := resolver.GetDeltasAhead(longerTree, lastCommonHash)
	if AheadErr != nil {
		fmt.Println("Deltas ahead error")
		return
	}

	deltasToPush = append(deltasToPush, deltasAhead...)

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
