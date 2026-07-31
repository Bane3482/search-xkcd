package core

import (
	searchpb "yadro.com/course/proto/search"
	updatepb "yadro.com/course/proto/update"
)

type Comics struct {
	ID    int64
	URL   string
	Words []string
}

func FromProtoComics(comics []*updatepb.ComicsInfo) []*Comics {
	result := make([]*Comics, 0)

	for _, c := range comics {
		result = append(result, &Comics{
			ID:    c.Id,
			URL:   c.Url,
			Words: c.Words,
		})
	}
	return result
}

func ToProtoComics(comics []*Comics) []*searchpb.Comics {
	result := make([]*searchpb.Comics, 0)

	for _, c := range comics {
		result = append(result, &searchpb.Comics{
			Id:  c.ID,
			Url: c.URL,
		})
	}
	return result
}
