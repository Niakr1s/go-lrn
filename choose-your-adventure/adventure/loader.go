package adventure

import (
	"os"
	"path"
	"strings"
)

type FsAdventureLoader struct {
	dirpath string
}

func NewFsAdventureLoader(dirpath string) *FsAdventureLoader {
	return &FsAdventureLoader{
		dirpath: dirpath,
	}
}

func (l *FsAdventureLoader) LoadAdventures() (Adventures, error) {
	adventures := Adventures{}

	dirEntries, err := os.ReadDir(l.dirpath)

	if err != nil {
		return nil, err
	}
	for _, item := range dirEntries {
		if item.IsDir() {
			continue
		}

		filename := item.Name()
		adventureName := strings.TrimRight(filename, ".json")
		if strings.HasSuffix(filename, ".json") {
			contents, err := os.ReadFile(path.Join(l.dirpath, filename))
			if err != nil {
				return nil, err
			}
			a, err := JsonToAdventure(contents)
			if err != nil {
				return nil, err
			}
			adventures[adventureName] = a

		}
	}
	return adventures, nil
}
