package utils

import "os"

const filePattern = "*.part"

func CreateTempFile(data []byte) (path string, err error) {
	f, err := os.CreateTemp("", filePattern)
	if err != nil {
		return
	}

	defer f.Close()

	path = f.Name()

	if _, err = f.Write(data); err != nil {
		return
	}

	return
}
