package server

import (
	"embed"
	"io"
	"lrn/choose-your-adventure/adventure"
	"text/template"
)

//go:embed templates
var templateFs embed.FS

type DefaultWriter struct {
	arcTempl   *template.Template
	indexTempl *template.Template
}

func NewDefaultWriter() (*DefaultWriter, error) {
	arcTempl, err := template.ParseFS(templateFs, "templates/arc.tmpl.html", "templates/head.tmpl.html")
	if err != nil {
		return nil, err
	}
	indexTempl, err := template.ParseFS(templateFs, "templates/index.tmpl.html", "templates/head.tmpl.html")
	if err != nil {
		return nil, err
	}
	res := &DefaultWriter{
		arcTempl:   arcTempl,
		indexTempl: indexTempl,
	}
	return res, nil
}

func (dw *DefaultWriter) WriteArc(w io.Writer, adventureName, arcName string, arc *adventure.Arc) error {
	type ArcWithData struct {
		*adventure.Arc
		AdventureName string
		ArcName       string
	}
	arcWithData := ArcWithData{
		Arc:           arc,
		AdventureName: adventureName,
		ArcName:       arcName,
	}
	err := dw.arcTempl.Execute(w, arcWithData)
	return err
}

func (dw *DefaultWriter) WriteIndex(w io.Writer, adventures adventure.Adventures) error {
	err := dw.indexTempl.Execute(w, adventures)
	return err
}
