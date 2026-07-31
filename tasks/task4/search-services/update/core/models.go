package core

import updatepb "yadro.com/course/proto/update"

type ServiceStatus string

const (
	StatusRunning ServiceStatus = "running"
	StatusIdle    ServiceStatus = "idle"
)

type DBStats struct {
	WordsTotal    int
	WordsUnique   int
	ComicsFetched int
}

type ServiceStats struct {
	DBStats
	ComicsTotal int
}

type Comics struct {
	ID    int
	URL   string
	Words []string
}

type XKCDInfo struct {
	ID          int    `json:"num"`
	URL         string `json:"img"`
	Title       string `json:"title"`
	Description string `json:"transcript"`
}

func ToProtoStatus(status ServiceStatus) updatepb.Status {
	switch status {
	case StatusIdle:
		return updatepb.Status_STATUS_IDLE
	case StatusRunning:
		return updatepb.Status_STATUS_RUNNING
	default:
		return updatepb.Status_STATUS_UNSPECIFIED
	}
}

func ToProtoStats(stats ServiceStats) *updatepb.StatsReply {
	return &updatepb.StatsReply{
		WordsTotal:    int64(stats.WordsTotal),
		WordsUnique:   int64(stats.WordsUnique),
		ComicsTotal:   int64(stats.ComicsTotal),
		ComicsFetched: int64(stats.ComicsFetched),
	}
}
