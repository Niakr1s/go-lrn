package pages

import (
	"fmt"
	"lrn/choose-your-adventure/cli/page"
)

type MockPage struct {
	LastWord string
}

func (p *MockPage) Display() {
	fmt.Printf("hello, your last word: %s\n", p.LastWord)
}

func (p *MockPage) HandleInput(input string) page.Page {
	p.LastWord = input
	if input == "r" {
		return &MockPage{}
	}
	return nil
}

func (p *MockPage) GetHelpMessage() string {
	return ""
}
