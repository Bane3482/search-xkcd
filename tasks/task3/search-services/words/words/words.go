package words

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	_ "embed"

	"github.com/kljensen/snowball"
)

//go:embed words.txt
var stopList string

type normalizer struct {
	StopMap map[string]struct{}
}

const (
	symbs = "!@\"#№$;%:^?&*()-_+=~`[]{}'\\|<>/., \n\r\t"
)

func New() *normalizer {
	stopMap := make(map[string]struct{})

	words := strings.Fields(stopList)

	for _, word := range words {
		stopMap[word] = struct{}{}
	}
	return &normalizer{
		StopMap: stopMap,
	}
}

func (n *normalizer) Norm(_ context.Context, phrase string) []string {
	if phrase == "" {
		return make([]string, 0)
	}

	words := strings.FieldsFunc(phrase, func(c rune) bool {
		return strings.Contains(symbs, string(c))
	})

	list := make(map[string]struct{})

	for _, word := range words {
		fmt.Println(word)
		stemmed, err := snowball.Stem(word, "english", true)

		if err != nil {
			return nil
		}

		if _, ok := n.StopMap[stemmed]; !ok {
			list[stemmed] = struct{}{}
		}
	}

	return slices.Collect(maps.Keys(list))
}
