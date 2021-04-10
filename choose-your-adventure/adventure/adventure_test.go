package adventure_test

import (
	"embed"
	"log"
	"lrn/choose-your-adventure/adventure"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:embed test_data
var testDataFs embed.FS

const testDatafolderName = "test_data"

type file struct {
	Contents []byte
	Name     string
}

func readFiles(foldername string) []file {
	dir, err := testDataFs.ReadDir(path.Join(testDatafolderName, foldername))
	if err != nil {
		log.Fatalf("can't read dir")
	}
	res := make([]file, 0)
	for _, f := range dir {
		fullName := path.Join(testDatafolderName, foldername, f.Name())
		contents, err := testDataFs.ReadFile(fullName)
		if err != nil {
			log.Fatalf("couldn't read file %s", fullName)
		}
		res = append(res, file{Contents: contents, Name: f.Name()})
	}
	return res
}

func readValidFiles() []file {
	return readFiles("valid")
}

func readInValidFiles() []file {
	return readFiles("invalid")
}

func TestJsonToAdventure(t *testing.T) {
	t.Run("valid adventures", func(t *testing.T) {
		for _, f := range readValidFiles() {
			t.Run(f.Name, func(t *testing.T) {
				a, err := adventure.JsonToAdventure(f.Contents)
				assert.Nil(t, err)
				assert.NotNil(t, a)
				assert.Equal(t, adventure.StartArc, a.CurrentArcName)

				assert.NotNil(t, a.GetCurrentArc())
			})
		}

	})
	t.Run("invalid adventures", func(t *testing.T) {
		for _, f := range readInValidFiles() {
			t.Run(f.Name, func(t *testing.T) {
				a, err := adventure.JsonToAdventure(f.Contents)
				assert.NotNil(t, err)
				assert.Nil(t, a)
			})
		}

	})
}
