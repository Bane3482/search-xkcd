package petname

import (
	petname "github.com/dustinkirkland/golang-petname"
)

func Generate(words int64, separator string) string {
	return petname.Generate(int(words), separator)
}

func GenerateMany(words int64, separator string, names int64) []string {
	list := make([]string, 0, names)

	for count := int64(0); count < names; count++ {
		name := petname.Generate(int(words), separator)

		list = append(list, name)
	}

	return list
}
