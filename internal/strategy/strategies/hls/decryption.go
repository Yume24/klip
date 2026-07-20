package hls

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strings"

	"github.com/Eyevinn/hls-m3u8/m3u8"
)

const noneMethod = "none"
const aes128Method = "aes-128"
const keyLength = 16
const ivLength = 16
const hexPrefix = "0x"

var errNoEncryption = errors.New("no encrpytion scheme")
var errNoAesEncryption = errors.New("no AES-128 encryption scheme")
var errInvalidIVLength = errors.New("invalid IV length")
var errInvalidKeyLength = errors.New("invalid key length")
var errInvalidBlockSize = errors.New("invalid block size")
var errBadPadding = errors.New("bad padding")

type decrpytionInfo struct {
	key []byte
	iv  []byte
}

func getAllKeys(playlist *m3u8.MediaPlaylist, playlistURL string) (map[int]decrpytionInfo, error) {
	playlistLength := playlist.Count()
	keyMap := make(map[int]decrpytionInfo, playlistLength)

	for i, segment := range playlist.GetAllSegments() {
		key, iv, err := getKeyFromSegment(segment, playlistURL)
		if err != nil {
			if errors.Is(err, errNoEncryption) {
				if i != 0 {
					keyMap[i] = decrpytionInfo{key: keyMap[i-1].key, iv: iv[:]}
				}
				continue
			} else {
				return nil, err
			}
		}

		keyMap[i] = decrpytionInfo{key: key[:], iv: iv[:]}
	}

	return keyMap, nil
}

func getKeyFromSegment(segment *m3u8.MediaSegment, playlistURL string) ([keyLength]byte, [ivLength]byte, error) {
	var resultKey [keyLength]byte
	var resultIV [ivLength]byte

	key, err := getAesEncryptionScheme(segment.Keys)
	if err != nil {
		return resultKey, resultIV, err
	}

	if key.IV != "" {
		resultIV, err = decodeIV(key.IV)
		if err != nil {
			return resultKey, resultIV, err
		}

	} else {
		binary.BigEndian.PutUint64(resultIV[ivLength/2:], segment.SeqId)
	}

	keyURI, err := resolveAbsoluteURL(playlistURL, key.URI)
	if err != nil {
		return resultKey, resultIV, err
	}

	keyBuf := &bytes.Buffer{}
	if err := getResponseBody(keyURI, keyBuf); err != nil {
		return resultKey, resultIV, err
	}

	if keyBuf.Len() != keyLength {
		return resultKey, resultIV, errInvalidKeyLength
	}

	resultKey = [keyLength]byte(keyBuf.Bytes())

	return resultKey, resultIV, nil
}

func decodeIV(iv string) ([ivLength]byte, error) {
	var result [ivLength]byte

	iv = strings.TrimPrefix(iv, hexPrefix)
	ivDecoded, err := hex.DecodeString(iv)
	if err != nil {
		return result, err
	}

	if len(ivDecoded) != ivLength {
		return result, errInvalidIVLength
	}

	result = [ivLength]byte(ivDecoded)
	return result, nil
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

func decryptSegment(segmentData, key, iv []byte) ([]byte, error) {
	if key == nil {
		return segmentData, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return segmentData, err
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	if len(segmentData)%aes.BlockSize != 0 {
		return segmentData, errInvalidBlockSize
	}

	mode.CryptBlocks(segmentData, segmentData)

	n := int(segmentData[len(segmentData)-1])
	if n < 1 || n > aes.BlockSize || n > len(segmentData) {
		return nil, errBadPadding
	}

	segmentData = segmentData[:len(segmentData)-n]

	return segmentData, nil
}
