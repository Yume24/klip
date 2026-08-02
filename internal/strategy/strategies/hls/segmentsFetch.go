package hls

import (
	"github.com/Eyevinn/hls-m3u8/m3u8"
	"github.com/Yume24/klip/internal/utils"
)

func getAllSegments(playlist *m3u8.MediaPlaylist, playlistURL string) ([]string, error) {
	keys, err := getAllKeys(playlist, playlistURL)
	if err != nil {
		return nil, err
	}
	jobs := buildFetchPlan(playlist, playlistURL, keys)
	return utils.RunFetchJobs(jobs)
}

func buildFetchPlan(playlist *m3u8.MediaPlaylist, playlistURL string, keys map[int]decrpytionInfo) []utils.FetchJob[string] {
	hasMap := playlist.Map != nil
	count := playlist.Count()

	if hasMap {
		count++
	}

	jobs := make([]utils.FetchJob[string], 0, count)

	if hasMap {
		jobs = append(jobs, func() (string, error) {
			return downloadMap(playlist.Map.URI, playlistURL)
		})
	}

	for i, segment := range playlist.GetAllSegments() {
		jobs = append(jobs, func() (string, error) {
			return downloadSegment(segment, playlistURL, keys[i])
		})
	}

	return jobs
}

func downloadSegment(segment *m3u8.MediaSegment, playlistURL string, decryption decrpytionInfo) (path string, err error) {
	segmentBuf, err := utils.ResolveURLAndDownload(playlistURL, segment.URI)
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

func downloadMap(mapURL, playlistURL string) (path string, err error) {
	mapData, err := utils.ResolveURLAndDownload(playlistURL, mapURL)
	if err != nil {
		return
	}

	path, err = utils.CreateTempFile(mapData.Bytes())
	if err != nil {
		return
	}

	return
}
