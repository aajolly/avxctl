/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
// package main

// import (
// 	"github.com/aajolly/avxctl/cmd"
// )

// func main() {
// 	cmd.Execute()
// }

package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed terraform/**
var embeddedFiles embed.FS

func MountEmbeddedFolderToTempDir() (string, error) {
	// Specify the desired name for the temporary folder
	tempFolderName := "embedded-terraform"

	// Get the default temporary directory
	tempDir := os.TempDir()

	// Append the desired folder name to the temporary directory
	tempFolderPath := filepath.Join(tempDir, tempFolderName)

	// Create a temporary directory
	// tempDir := "/tmp/embedded-terraform"
	err := os.Mkdir(tempFolderPath, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		// Remove the temporary directory if it already exists
		fmt.Printf("Removing directory: %s", tempFolderPath)
		os.RemoveAll(tempFolderPath)
	}

	// Create the temporary directory again
	fmt.Printf("Recreating directory: %s", tempFolderPath)
	err = os.MkdirAll(tempFolderPath, os.ModePerm)
	if err != nil {
		fmt.Println("2")
		fmt.Println(err)
		return "", err
	}

	// Walk through the embedded folder and copy all folders and files
	err = fs.WalkDir(embeddedFiles, "terraform", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("3")
			fmt.Println(err)
			return err
		}

		relativePath, err := filepath.Rel("terraform", path)
		if err != nil {
			fmt.Println("4")
			fmt.Println(err)
			return err
		}

		// Check if the entry is a directory
		if d.IsDir() {
			err = os.MkdirAll(filepath.Join(tempFolderPath, relativePath), os.ModePerm)
			if err != nil {
				fmt.Println("5")
				fmt.Println(err)
				return err
			}
		} else {
			// Copy the file from the embedded folder to the temporary directory
			src, err := embeddedFiles.Open(path)
			if err != nil {
				fmt.Println("6")
				fmt.Println(err)
				return err
			}
			defer src.Close()

			dst, err := os.Create(filepath.Join(tempFolderPath, relativePath))
			if err != nil {
				fmt.Println("7")
				fmt.Println(err)
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				fmt.Println("8")
				fmt.Println(err)
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Println("9")
		fmt.Println(err)
		return "", err
	}

	return tempFolderPath, nil
}

func main() {
	tempDir, err := MountEmbeddedFolderToTempDir()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Temporary folder created:", tempDir)
	// Clean up the temporary folder when you're done
	// defer os.RemoveAll(tempDir)
}
