package processor

import (
	"crypto/sha256"
	"fmt"
	"io"
)

func ComputeHash(fileReader io.Reader) (string, error) {
	hash := sha256.New()
	if _, err := io.Copy(hash, fileReader); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
