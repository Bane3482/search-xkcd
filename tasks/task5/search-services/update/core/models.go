package core

import (
	"github.com/lib/pq"
	updatepb "yadro.com/course/proto/update"
)

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

type ComicsPQ struct {
	ID    int            `db:"comics_id"`
	URL   string         `db:"comics_url"`
	Words pq.StringArray `db:"words"`
}

type Comics struct {
	ID    int
	URL   string
	Words []string
}

type XKCDInfo struct {
	ID         int    `json:"num"`
	URL        string `json:"img"`
	Title      string `json:"title"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
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

func ToProtoComics(comics []Comics) []*updatepb.ComicsInfo {
	result := make([]*updatepb.ComicsInfo, 0)

	for _, c := range comics {
		result = append(result, &updatepb.ComicsInfo{
			Id:    int64(c.ID),
			Url:   c.URL,
			Words: c.Words,
		})
	}

	return result
}

func FromPQArray(comics []ComicsPQ) []Comics {
	result := make([]Comics, len(comics))
	for i, c := range comics {
		result[i] = Comics{
			ID:    c.ID,
			URL:   c.URL,
			Words: []string(c.Words),
		}
	}

	return result
}
