package core

import updatepb "yadro.com/course/proto/update"

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
