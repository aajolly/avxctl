package terraformhelper

import (
	"embed"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

//go:embed modules/**
var embeddedFiles embed.FS

func MountEmbeddedFolderToTempDir() (string, error) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "embedded-terraform")
	if err != nil {
		return "", err
	}

	// Remove the temporary directory if it already exists
	err = os.RemoveAll(tempDir)
	if err != nil {
		return "", err
	}

	// Create the temporary directory again
	err = os.MkdirAll(tempDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Walk through the embedded folder and copy all folders and files
	err = fs.WalkDir(embeddedFiles, "modules", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relativePath, err := filepath.Rel("modules", path)
		if err != nil {
			return err
		}

		// Check if the entry is a directory
		if d.IsDir() {
			err = os.MkdirAll(filepath.Join(tempDir, relativePath), os.ModePerm)
			if err != nil {
				return err
			}
		} else {
			// Copy the file from the embedded folder to the temporary directory
			src, err := embeddedFiles.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()

			dst, err := os.Create(filepath.Join(tempDir, relativePath))
			if err != nil {
				return err
			}
			defer dst.Close()

			_, err = io.Copy(dst, src)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	return tempDir, nil
}
