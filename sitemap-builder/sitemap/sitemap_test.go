package sitemap_test

import (
	"context"
	"lrn/sitemap-builder/sitemap"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_GenerateSitemap(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*1000)
	defer cancel()

	res, err := sitemap.Build(ctx, "https://yandex.ru")

	t.Logf("%d links", len(res.Links))

	assert.NoError(t, err)
	assert.NotNil(t, res)
}
