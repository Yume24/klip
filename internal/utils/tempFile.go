package utils

import "os"

const filePattern = "*.part"

func CreateTempFile(data []byte) (string, error) {
	var path string

	f, err := os.CreateTemp("", filePattern)
	if err != nil {
		return path, err
	}

	defer f.Close()

	path = f.Name()

	if _, err := f.Write(data); err != nil {
		return path, err
	}

	return path, nil
}
