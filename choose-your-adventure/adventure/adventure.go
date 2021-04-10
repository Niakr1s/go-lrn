package adventure

import (
	"encoding/json"
	"fmt"
)

type Option struct {
	Text    string `json:"text,"`
	ArcName string `json:"arc"`
}

type Arc struct {
	Title   string    `json:"title"`
	Story   []string  `json:"story"`
	Options []*Option `json:"options"`
}

type Arcs map[string]*Arc

func (a Arcs) FindArc(arcName string) (*Arc, error) {
	arc, ok := a[arcName]
	if !ok {
		return nil, fmt.Errorf("arc with name %s not found", arcName)
	}
	return arc, nil
}

func (s Arcs) checkArc(arcName string) error {
	arc, ok := s[arcName]
	if !ok {
		return fmt.Errorf("no arc with name %s", arcName)
	}

	if arc.Title == "" {
		return fmt.Errorf("arc %s has no title", arcName)
	}

	if len(arc.Story) == 0 {
		return fmt.Errorf("arc %s has no contents", arcName)
	}

	for i, option := range arc.Options {
		if option.Text == "" {
			return fmt.Errorf("%d option of arc %s has no text", i, arcName)
		}
		if option.ArcName == "" {
			return fmt.Errorf("%d option of arc %s has no arc", i, arcName)
		}
		nextArcName := option.ArcName
		if nextArcName == arcName {
			continue
		}
		err := s.checkArc(nextArcName)
		if err != nil {
			return err
		}
	}
	return nil
}

const StartArc = "intro"

type Adventure struct {
	Arcs           Arcs
	CurrentArcName string
}

type Adventures map[string]*Adventure

func (a Adventures) FindArc(advName string, arcName string) (*Arc, error) {
	adv, ok := a[advName]
	if !ok {
		return nil, fmt.Errorf("adventure with name %s not found", advName)
	}
	arc, err := adv.Arcs.FindArc(arcName)
	if !ok {
		return nil, err
	}
	return arc, nil
}

func newAdventure(stories Arcs) (*Adventure, error) {
	a := &Adventure{
		Arcs:           stories,
		CurrentArcName: StartArc,
	}
	err := a.check()
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *Adventure) check() error {
	return a.Arcs.checkArc(a.CurrentArcName)
}

func (a *Adventure) GetCurrentArc() *Arc {
	return a.Arcs[a.CurrentArcName]
}

func JsonToAdventure(j []byte) (*Adventure, error) {
	stories := make(Arcs)
	err := json.Unmarshal(j, &stories)
	if err != nil {
		return nil, err
	}
	adventure, err := newAdventure(stories)
	if err != nil {
		return nil, err
	}
	return adventure, nil
}
