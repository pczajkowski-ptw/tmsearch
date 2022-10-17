package main

import (
	"sort"
)

// Language contains information about language
type Language struct {
	Name string
	Code string
}

func buildArrayOfLanguages(languagesMap map[string]string) []Language {
	var array []Language

	for key, value := range languagesMap {
		array = append(array, Language{Name: value, Code: key})
	}

	return array
}

func SortLanguagesByName(array []Language) {
	sort.Slice(array, func(i, j int) bool { return array[i].Name < array[j].Name })
}
