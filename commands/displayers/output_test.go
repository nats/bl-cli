package displayers

import (
	"bytes"
	"testing"

	"github.com/binarylane/bl-cli/bl"

	"github.com/stretchr/testify/assert"
)

func TestDisplayerDisplay(t *testing.T) {
	emptyImages := make([]bl.Image, 0)
	var nilImages []bl.Image

	tests := []struct {
		name         string
		item         Displayable
		expectedJSON string
	}{
		{
			name:         "displaying a non-nil slice of Volumes should return an empty JSON array",
			item:         &Image{Images: emptyImages},
			expectedJSON: `[]`,
		},
		{
			name:         "displaying a nil slice of Volumes should return an empty JSON array",
			item:         &Image{Images: nilImages},
			expectedJSON: `[]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}

			displayer := Displayer{
				OutputType: "json",
				Item:       tt.item,
				Out:        out,
			}

			err := displayer.Display()
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedJSON, out.String())
		})
	}
}
