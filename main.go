package main

import (
	"fmt"
	"io"
	"os"

	"filippo.io/age"
	shell "github.com/ipfs/go-ipfs-api"
)

func main() {
    // Input file to be read
    file, err := os.Open("file.txt")
    // Output encrypted file
    f, _ := os.Create("out.age")
    identity, err := age.GenerateX25519Identity()
    if err != nil {
        fmt.Println("Failed to generate key")
    }
    publicKey := identity.Recipient().String()
    fmt.Printf("Public key: %s...\n",publicKey)
    // recipient, err := age.ParseX25519Recipient(publicKey)
    // Encrypt file and set for my key to be allowed for decryption
    w, err := age.Encrypt(f, identity.Recipient())
    if err != nil{
        fmt.Println("Failed to create file")
    }
    if _, err := io.Copy(w, file); err != nil{
        fmt.Println("Unable to copy")
    }
    file.Close()
    if err := w.Close(); err != nil{
        fmt.Println("failed to close")
    }

    // Add to IPFS node
    sh := shell.NewShell("localhost:5001")
    cid, err := sh.Add(file)
    if err != nil {
        fmt.Fprintf(os.Stderr, "error: %s", err)
        os.Exit(1)
    }
    fmt.Printf("%s", cid)
    
}
