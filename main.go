package main

import (
	"context"
	"filippo.io/age"
	"flag"
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mholt/archiver/v4"
	"io"
	"os"
)

func generate_tar(dir_name, output_name string) {
	files, err := archiver.FilesFromDisk(nil, map[string]string{
		dir_name: "",
	})
	if err != nil {
		fmt.Println(err)
	}
	out, err := os.Create(output_name)
	if err != nil {
		fmt.Println(err)
	}
	defer out.Close()
	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}
	err = format.Archive(context.Background(), out, files)
}

func encrypt_file(input_name, output_name, keyfile string) bool {
	// Input file to be read
	file, _ := os.Open(input_name)
	// Output encrypted file
	f, _ := os.Create(output_name)
	keyFile, _ := os.Open(keyfile)
	identities, err := age.ParseRecipients(keyFile)
	// Encrypt file and set for keys in my
	w, err := age.Encrypt(f, identities...)
	if err != nil {
		fmt.Println("Failed to create file")
	}
	if _, err := io.Copy(w, file); err != nil {
		fmt.Println("Unable to copy")
	}
	defer file.Close()
	defer keyFile.Close()
	if err := w.Close(); err != nil {
		fmt.Println("failed to close")
		return false
	}
	return true

}

func main() {
	key_in := flag.String("k", "key.txt", "key file")
	dir := flag.String("d", "./", "input directory")
	flag.Parse()
	temp_name := "out.age"
	in_name := "out.tar.gz"
	generate_tar(*dir, in_name)
	encrypt_file(in_name, temp_name, *key_in)
	file, err := os.Open(temp_name)
	// Add to IPFS node
	sh := shell.NewShell("localhost:5001")
	cid, err := sh.Add(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("IPFS Hash: %s\n", cid)
	file.Close()
}
