package sitemap_test

import (
	"context"
	"lrn/sitemap-builder/sitemap"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateSitemap(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	res, err := sitemap.Build(ctx, "https://www.calhoun.io")

	t.Logf("%d links", len(res.Links))

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
