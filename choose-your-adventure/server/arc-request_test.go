package server_test

import (
	"fmt"
	"lrn/choose-your-adventure/adventure"
	"lrn/choose-your-adventure/server"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArcRequestFromUrl(t *testing.T) {
	cases := []struct {
		url       string
		wantError bool
		expected  server.ArcRequest
	}{
		{"/", true, server.ArcRequest{}},
		{"/name", false, server.ArcRequest{Name: "name", Arc: adventure.StartArc}},
		{"/name/", false, server.ArcRequest{Name: "name", Arc: adventure.StartArc}},
		{"/name/arc", false, server.ArcRequest{Name: "name", Arc: "arc"}},
		{"/name/arc/", false, server.ArcRequest{Name: "name", Arc: "arc"}},
		{"/name/arc/smth", false, server.ArcRequest{Name: "name", Arc: "arc"}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("url: %s", c.url), func(t *testing.T) {
			got, err := server.GetArcRequestFromUrl(c.url)
			if err != nil {
				assert.True(t, c.wantError)
				return
			}
			assert.Equal(t, c.expected, got)
		})
	}
}
