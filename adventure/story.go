package adventure

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type arcLink struct {
	Text string
	Arc  string
}

type storyArc struct {
	Title   string
	Story   []string
	Options []arcLink
}

func readFile(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func ReadStory(filePath string) map[string]storyArc {
	file := readFile(filePath)
	storyMap := make(map[string]storyArc)
	err := json.Unmarshal(file, &storyMap)
	if err != nil {
		log.Fatal(err)
	}
	return storyMap
}
