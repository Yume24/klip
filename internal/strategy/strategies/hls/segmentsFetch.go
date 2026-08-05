package hls

import (
	"bytes"

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

	if hasMap {
		url, err := utils.ResolveAbsoluteURL(playlistURL, playlist.Map.URI)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, func() (string, error) {
			return downloadMap(url, playlist.Map.Limit, playlist.Map.Offset)
		})
	}

	for i, segment := range playlist.GetAllSegments() {
		url, err := utils.ResolveAbsoluteURL(playlistURL, segment.URI)
		if err != nil {
			return nil, err
		}

		jobs = append(jobs, func() (string, error) {
			return downloadSegment(url, segment.Limit, segment.Offset, keys[i])
		})
	}

	return jobs, nil
}

func downloadSegment(segmentURL string, length, offset int64, decryption decrpytionInfo) (path string, err error) {
	var fetchOpts []utils.RequestOption
	segmentBuf := &bytes.Buffer{}

	if length > 0 {
		fetchOpts = append(fetchOpts, utils.WithByteRange(utils.ByteRange{Offset: offset, Length: length}))
	}

	err = utils.GetResponseBody(segmentURL, segmentBuf, fetchOpts...)
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

func downloadMap(mapURL string, length, offset int64) (path string, err error) {
	var fetchOpts []utils.RequestOption
	mapBuf := &bytes.Buffer{}

	if length > 0 {
		fetchOpts = append(fetchOpts, utils.WithByteRange(utils.ByteRange{Offset: offset, Length: length}))
	}

	err = utils.GetResponseBody(mapURL, mapBuf, fetchOpts...)
	if err != nil {
		return
	}

	path, err = utils.CreateTempFile(mapBuf.Bytes())
	if err != nil {
		return
	}

	return
}
