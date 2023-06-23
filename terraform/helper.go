package terraformhelper

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

//go:embed modules/**
var embeddedFiles embed.FS

// Define DateTime for logging
var dt = time.Now().Format("01-02-2006 15:04:05")

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
		color.Red(dt + ": " + err.Error())

		// Remove the temporary directory if it already exists
		color.Blue(dt+": "+"Removing directory: %s", tempFolderPath)
		os.RemoveAll(tempFolderPath)
	}

	// Create the temporary directory again
	color.Blue("Recreating directory: %s", tempFolderPath)
	err = os.MkdirAll(tempFolderPath, os.ModePerm)
	if err != nil {
		color.Red(dt + ": " + err.Error())
		return "", err
	}

	// Walk through the embedded folder and copy all folders and files
	err = fs.WalkDir(embeddedFiles, "modules", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			color.Red(dt + ": " + err.Error())
			return err
		}

		relativePath, err := filepath.Rel("modules", path)
		if err != nil {
			color.Red(dt + ": " + err.Error())
			return err
		}

		// Check if the entry is a directory
		if d.IsDir() {
			err = os.MkdirAll(filepath.Join(tempFolderPath, relativePath), os.ModePerm)
			if err != nil {
				color.Red(dt + ": " + err.Error())
				return err
			}
		} else {
			// Copy the file from the embedded folder to the temporary directory
			src, err := embeddedFiles.Open(path)
			if err != nil {
				color.Red(dt + ": " + err.Error())
				return err
			}
			defer src.Close()

			dst, err := os.Create(filepath.Join(tempFolderPath, relativePath))
			if err != nil {
				color.Red(dt + ": " + err.Error())
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				color.Red(dt + ": " + err.Error())
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return tempFolderPath, nil
}
