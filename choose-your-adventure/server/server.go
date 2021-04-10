package server

import (
	"io"
	"log"
	"lrn/choose-your-adventure/adventure"
	"net/http"
)

type AdventureLoader interface {
	LoadAdventures() (adventure.Adventures, error)
}

type ArcWriter interface {
	WriteArc(w io.Writer, adventureName, arcName string, arc *adventure.Arc) error
}

type IndexWriter interface {
	WriteIndex(w io.Writer, adventures adventure.Adventures) error
}

type ServerOptions struct {
	AdventureLoader AdventureLoader
	ArcWriter       ArcWriter
	IndexWriter     IndexWriter
}

func (so *ServerOptions) normalize() {
	if so.AdventureLoader == nil {
		so.AdventureLoader = EmptyAdventureLoader{}
	}
	if so.ArcWriter == nil {
		defaultArcWriter, err := NewDefaultWriter()
		if err != nil {
			log.Fatalf("couldn't load defaultArcWriter: %s", err)
		}
		so.ArcWriter = defaultArcWriter
	}
	if so.IndexWriter == nil {
		defaultIndexWriter, err := NewDefaultWriter()
		if err != nil {
			log.Fatalf("couldn't load defaultArcWriter: %s", err)
		}
		so.IndexWriter = defaultIndexWriter
	}
}

type Server struct {
	adventures  adventure.Adventures
	arcWriter   ArcWriter
	indexWriter IndexWriter
}

func NewServer(serverOptions ServerOptions) (*Server, error) {
	serverOptions.normalize()
	a, err := serverOptions.AdventureLoader.LoadAdventures()
	if err != nil {
		return nil, err
	}

	s := &Server{
		adventures:  a,
		arcWriter:   serverOptions.ArcWriter,
		indexWriter: serverOptions.IndexWriter,
	}
	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		s.indexWriter.WriteIndex(w, s.adventures)
		return
	}

	arcRequest, err := GetArcRequestFromUrl(r.URL.Path)
	if err != nil {
		s.handleNotFound(w, r)
		return
	}

	arc, err := s.adventures.FindArc(arcRequest.Name, arcRequest.Arc)
	if err != nil {
		s.handleNotFound(w, r)
		return
	}

	err = s.arcWriter.WriteArc(w, arcRequest.Name, arcRequest.Arc, arc)
	if err != nil {
		log.Printf("error while writing arc: %v", err)
	}
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)

}
