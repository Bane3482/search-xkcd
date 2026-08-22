package core

import (
	searchpb "yadro.com/course/proto/search"
	updatepb "yadro.com/course/proto/update"
)

type StatusUpdate string

const (
	StatusUpdateUnknown StatusUpdate = "unknown"
	StatusUpdateIdle    StatusUpdate = "idle"
	StatusUpdateRunning StatusUpdate = "running"
)

type UpdateStats struct {
	WordsTotal    int `json:"words_total"`
	WordsUnique   int `json:"words_unique"`
	ComicsFetched int `json:"comics_fetched"`
	ComicsTotal   int `json:"comics_total"`
}

type UpdateStatus struct {
	Status StatusUpdate `json:"status"`
}

type Comics struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

type ComicsReply struct {
	Comics []Comics `json:"comics"`
	Total  int      `json:"total"`
}

type PingResponse struct {
	Replies map[string]string `json:"replies"`
}

func FromProtoStatus(status updatepb.Status) StatusUpdate {
	switch status {
	case updatepb.Status_STATUS_IDLE:
		return StatusUpdateIdle
	case updatepb.Status_STATUS_RUNNING:
		return StatusUpdateRunning
	default:
		return StatusUpdateUnknown
	}
}

func FromProtoComicsReply(reply *searchpb.SearchReply) ComicsReply {
	comics := make([]Comics, 0)

	for _, c := range reply.Comics {
		comics = append(comics, Comics{
			ID:  int(c.Id),
			URL: c.Url,
		})
	}
	return ComicsReply{
		Comics: comics,
		Total:  int(reply.Total),
	}
}
