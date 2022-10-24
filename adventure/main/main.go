package main

import (
	"adventure"
)

func main() {
	storyMap := adventure.ReadStory("../data/gopher.json")
	_ = storyMap
}
