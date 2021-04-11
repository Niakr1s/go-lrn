package pages

import (
	"fmt"
	"lrn/choose-your-adventure/adventure"
	"lrn/choose-your-adventure/cli/page"
)

type AdventuresIndexPage struct {
	Adventures adventure.Adventures
}

func NewAdventuresIndexPage(adv adventure.Adventures) *AdventuresIndexPage {
	return &AdventuresIndexPage{Adventures: adv}
}

func (p *AdventuresIndexPage) Display() {
	fmt.Println("List of adventures:")
	for advName := range p.Adventures {
		fmt.Println(advName)
	}
	fmt.Printf("\nEnter name of adventure to start: ")
}

func (p *AdventuresIndexPage) GetHelpMessage() string {
	return ""
}

func (p *AdventuresIndexPage) HandleInput(input string) page.Page {
	arcPage, err := NewAdventuresArcPage(p.Adventures, input, "")
	if err != nil {
		return nil
	}
	return arcPage
}

type AdventuresArcPage struct {
	Adventures    adventure.Adventures
	adventureName string
	arcName       string
	arc           *adventure.Arc
}

func NewAdventuresArcPage(adv adventure.Adventures, advName, arcName string) (*AdventuresArcPage, error) {
	if arcName == "" {
		arcName = adventure.StartArc
	}
	arc, err := adv.FindArc(advName, arcName)
	if err != nil {
		return nil, err
	}
	return &AdventuresArcPage{
		Adventures:    adv,
		adventureName: advName,
		arcName:       arcName,
		arc:           arc,
	}, nil
}

func (p *AdventuresArcPage) getIntroPage() *AdventuresArcPage {
	arc, _ := NewAdventuresArcPage(p.Adventures, p.adventureName, adventure.StartArc)
	return arc
}

func (p *AdventuresArcPage) getNextArcPage(arcName string) (*AdventuresArcPage, error) {
	if !p.arc.HasOption(arcName) {
		return nil, fmt.Errorf("arc has no such option")
	}

	arc, err := NewAdventuresArcPage(p.Adventures, p.adventureName, arcName)
	if err != nil {
		return nil, err
	}
	return arc, nil
}

func (p *AdventuresArcPage) Display() {
	fmt.Printf("Title: %s\n", p.arc.Title)

	fmt.Println("\nStory:")
	for _, storyLine := range p.arc.Story {
		fmt.Println(storyLine)
	}

	if len(p.arc.Options) > 0 {
		fmt.Println("\nOptions:")
		for _, opt := range p.arc.Options {
			fmt.Printf("%s: %s\n\n", opt.ArcName, opt.Text)
		}
		fmt.Printf("Enter name of option to continue: ")
	} else {
		fmt.Printf("\nEnter intro to load index page: ")
	}
}

func (p *AdventuresArcPage) GetHelpMessage() string {
	return "!intro: return to intro page"
}

func (p *AdventuresArcPage) HandleInput(input string) page.Page {
	switch input {
	case "!intro":
		return p.getIntroPage()

	default:
		if nextPage, err := p.getNextArcPage(input); err == nil {
			return nextPage
		}
	}

	return nil
}
