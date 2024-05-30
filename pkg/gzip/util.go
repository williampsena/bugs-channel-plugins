package gzip

import (
	"compress/gzip"
	"io"
)

// Decompress content compressed with the gzip format.
func UnzipReader(reader io.Reader) (string, error) {
	zipReader, err := gzip.NewReader(reader)

	if err != nil {
		return "", err
	}

	output, err := io.ReadAll(zipReader)

	if err != nil {
		return "", err
	}

	return string(output), nil
}
