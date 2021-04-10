package adventure_test

import (
	"lrn/choose-your-adventure/adventure"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const adventuresDirEnvKey = "ADVENTURES_DIR"

func TestFsLoader(t *testing.T) {
	dirpath := os.Getenv(adventuresDirEnvKey)
	if dirpath == "" {
		t.Fatalf("for test FsAdventureLoader please provide %s env variable with gopher.json adventure inside it", adventuresDirEnvKey)
		return
	}
	a, err := adventure.NewFsAdventureLoader(dirpath).LoadAdventures()
	assert.NoError(t, err)
	_, hasGopher := a["gopher"]
	assert.True(t, hasGopher)
}
