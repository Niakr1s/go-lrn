package server

import "lrn/choose-your-adventure/adventure"

type EmptyAdventureLoader struct{}

func (eal EmptyAdventureLoader) LoadAdventures() (adventure.Adventures, error) {
	return adventure.Adventures{}, nil
}
