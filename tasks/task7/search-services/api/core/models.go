package core

import (
	searchpb "yadro.com/course/proto/search"
	updatepb "yadro.com/course/proto/update"
)

type UpdateStatus string

const (
	StatusUpdateUnknown UpdateStatus = "unknown"
	StatusUpdateIdle    UpdateStatus = "idle"
	StatusUpdateRunning UpdateStatus = "running"
)

type UpdateStats struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func FromProtoStatus(status updatepb.Status) UpdateStatus {
	switch status {
	case updatepb.Status_STATUS_IDLE:
		return StatusUpdateIdle
	case updatepb.Status_STATUS_RUNNING:
		return StatusUpdateRunning
	default:
		return StatusUpdateUnknown
	}
}

func FromProtoComicsReply(reply *searchpb.SearchReply) []Comics {
	comics := make([]Comics, 0)

	for _, c := range reply.Comics {
		comics = append(comics, Comics{
			ID:  int(c.Id),
			URL: c.Url,
		})
	}
	return comics
}
