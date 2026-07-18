package hls

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"fmt"
	"strings"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const noneMethod = "none"
const aes128Method = "aes-128"
const ivLength = 16

var errNoEncryption = errors.New("no encrpytion scheme")
var errNoAesEncryption = errors.New("no AES-128 encryption scheme")
var invalidIVLength = errors.New("invalid IV length")

func decryptSegment(segment *orderedMediaSegment) error {
	if len(segment.iv) != ivLength {
		return invalidIVLength
	}

	if segment.key == nil {
		return nil
	}

	block, err := aes.NewCipher(segment.key)
	if err != nil {
		return err
	}

	mode := cipher.NewCBCDecrypter(block, segment.iv)
	mode.CryptBlocks(segment.data, segment.data)

	n := int(segment.data[len(segment.data)-1])
	if n < 1 || n > aes.BlockSize || n > len(segment.data) {
		return fmt.Errorf("segment %d: bad pad n=%d len=%d", segment.id, n, len(segment.data))
	}
	segment.data = segment.data[:len(segment.data)-n]

	return nil
}

func getAesEncryptionScheme(keys []m3u8.Key) (m3u8.Key, error) {
	var foundKey m3u8.Key

	if len(keys) == 0 {
		return foundKey, errNoEncryption
	}

	for _, key := range keys {
		switch strings.ToLower(key.Method) {
		case aes128Method:
			foundKey = key
			return foundKey, nil
		case noneMethod:
			return foundKey, errNoEncryption
		}
	}

	return foundKey, errNoAesEncryption
}
