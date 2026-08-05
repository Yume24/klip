package hls

import (
	"bytes"
	"fmt"
	"math"

	"github.com/Eyevinn/hls-m3u8/m3u8"
	"github.com/Yume24/klip/internal/utils"
)

func getAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) ([]string, error) {
	keys, err := getAllKeys(playlist, playlistURL)
	if err != nil {
		return nil, err
	}
	jobs, err := buildFetchPlan(playlist, playlistURL, keys)
	if err != nil {
		return nil, err
	}
	return utils.RunFetchJobs(jobs)
}

func buildFetchPlan(playlist *m3u8.MediaPlaylist, playlistURL string, keys map[int]decrpytionInfo) ([]utils.FetchJob[string], error) {
	hasMap := playlist.Map != nil
	count := playlist.Count()

	if hasMap {
		count++
	}

	jobs := make([]utils.FetchJob[string], 0, count)
	offsets := make(map[string]int64, count)

	if hasMap {
		url, err := utils.ResolveAbsoluteURL(playlistURL, playlist.Map.URI)
		if err != nil {
			return nil, err
		}

		byteRange, err := resolveRange(url, playlist.Map.Limit, playlist.Map.Offset, offsets)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, func() (string, error) {
			return downloadMap(url, byteRange)
		})
	}

	for i, segment := range playlist.GetAllSegments() {
		url, err := utils.ResolveAbsoluteURL(playlistURL, segment.URI)
		if err != nil {
			return nil, err
		}

		byteRange, err := resolveRange(url, playlist.Map.Limit, playlist.Map.Offset, offsets)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, func() (string, error) {
			return downloadSegment(url, byteRange, keys[i])
		})
	}

	return jobs, nil
}

func resolveRange(url string, limit, offset int64, offsets map[string]int64) (br utils.ByteRange, err error) {
	if limit <= 0 {
		return // Fetch whole resource
	}

	if offset < 0 || offset > math.MaxInt64-limit {
		return br, fmt.Errorf("byterange out of range: %d@%d", limit, offset)
	}

	curr, seen := offsets[url]
	if seen && offset == 0 {
		br.Offset = curr
	}

	if offset < curr {
		return br, fmt.Errorf("byterange %d@%d overlaps prior range in %s", limit, offset, url)
	}

	offsets[url] = offset + limit
	br.Length = limit

	return
}

func downloadSegment(segmentURL string, br utils.ByteRange, decryption decrpytionInfo) (path string, err error) {
	segmentBuf := &bytes.Buffer{}

	err = utils.GetResponseBody(segmentURL, segmentBuf, utils.WithByteRange(br))
	if err != nil {
		return
	}

	decryptedSegment, err := decryptSegment(segmentBuf.Bytes(), decryption.key, decryption.iv)
	if err != nil {
		return
	}

	path, err = utils.CreateTempFile(decryptedSegment)
	if err != nil {
		return
	}

	return
}

func downloadMap(mapURL string, br utils.ByteRange) (path string, err error) {
	mapBuf := &bytes.Buffer{}

	err = utils.GetResponseBody(mapURL, mapBuf, utils.WithByteRange(br))
	if err != nil {
		return
	}

	path, err = utils.CreateTempFile(mapBuf.Bytes())
	if err != nil {
		return
	}

	return
}
