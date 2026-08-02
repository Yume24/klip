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
	"github.com/Yume24/klip/internal/utils"
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

	var currentKeyData []byte
	var currentKey m3u8.Key
	var currentIV []byte

	for i, segment := range playlist.GetAllSegments() {
		segmentKey, err := getAesEncryptionScheme(segment.Keys)
		if err != nil && !errors.Is(err, errNoEncryption) {
			return nil, err
		} else if err == nil {
			currentKey = segmentKey
			currentIV = nil

			keyData, err := deriveKeyData(currentKey, playlistURL)
			if err != nil {
				return nil, err
			}

			currentKeyData = keyData[:]

		}

		if currentKey.IV == "" {
			iv := deriveIVFromSegment(segment)
			currentIV = iv[:]
		} else if currentIV == nil {
			iv, err := deriveIVFromKey(currentKey)
			if err != nil {
				return nil, err
			}
			currentIV = iv[:]
		}

		keyMap[i] = decrpytionInfo{key: currentKeyData, iv: currentIV}
	}

	return keyMap, nil
}

func deriveIVFromKey(key m3u8.Key) (resultIV [ivLength]byte, err error) {
	iv := strings.TrimPrefix(key.IV, hexPrefix)
	ivDecoded, err := hex.DecodeString(iv)
	if err != nil {
		return
	}

	if len(ivDecoded) != ivLength {
		return resultIV, errInvalidIVLength
	}

	resultIV = [ivLength]byte(ivDecoded)
	return
}

func deriveIVFromSegment(segment *m3u8.MediaSegment) (resultIV [ivLength]byte) {
	binary.BigEndian.PutUint64(resultIV[ivLength/2:], segment.SeqId)
	return
}

func deriveKeyData(key m3u8.Key, playlistURL string) (resultKey [keyLength]byte, err error) {
	keyURI, err := utils.ResolveAbsoluteURL(playlistURL, key.URI)
	if err != nil {
		return
	}

	keyBuf := &bytes.Buffer{}
	if err = utils.GetResponseBody(keyURI, keyBuf); err != nil {
		return
	}

	if keyBuf.Len() != keyLength {
		return resultKey, errInvalidKeyLength
	}

	resultKey = [keyLength]byte(keyBuf.Bytes())

	return
}

func getAesEncryptionScheme(keys []m3u8.Key) (foundKey m3u8.Key, _ error) {
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
