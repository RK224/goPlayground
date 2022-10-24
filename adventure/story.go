package adventure

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type ArcLink struct {
	Text string
	Arc  string
}

type StoryArc struct {
	Title   string
	Story   []string
	Options []ArcLink
}

func readFile(filePath string) []byte {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func ReadStory(filePath string) map[string]StoryArc {
	file := readFile(filePath)
	storyMap := make(map[string]StoryArc)
	err := json.Unmarshal(file, &storyMap)
	if err != nil {
		log.Fatal(err)
	}
	return storyMap
}
