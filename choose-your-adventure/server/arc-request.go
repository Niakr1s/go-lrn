package server

import (
	"fmt"
	"lrn/choose-your-adventure/adventure"
	"strings"
)

type ArcRequest struct {
	Name string
	Arc  string
}

func (ar ArcRequest) isValid() bool {
	return ar.Name != "" && ar.Arc != ""
}

func GetArcRequestFromUrl(url string) (ArcRequest, error) {
	res := ArcRequest{
		Name: "",
		Arc:  adventure.StartArc,
	}
	splitted := strings.Split(strings.Trim(url, "/"), "/")
	if len(splitted) == 0 {
		return ArcRequest{}, fmt.Errorf("invalid url")
	}
	res.Name = splitted[0]
	if len(splitted) > 1 {
		res.Arc = splitted[1]
	}
	if !res.isValid() {
		return ArcRequest{}, fmt.Errorf("invalid url")
	}
	return res, nil
}
