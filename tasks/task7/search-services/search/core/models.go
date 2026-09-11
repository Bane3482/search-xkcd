package core

import (
	searchpb "yadro.com/course/proto/search"
)

type Comics struct {
	ID    int64
	URL   string
	Words []string
}

type ComicsInfo struct {
	ID  int64  `db:"comics_id"`
	URL string `db:"comics_url"`
}

func ToProtoComics(comics []ComicsInfo) []*searchpb.Comics {
	result := make([]*searchpb.Comics, 0)

	for _, c := range comics {
		result = append(result, &searchpb.Comics{
			Id:  c.ID,
			Url: c.URL,
		})
	}
	return result
}
