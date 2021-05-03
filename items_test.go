package main

import (
	"encoding/json"
	"io/fs"
	"os"
	"testing"

	"github.com/creativeprojects/gopenHAB/api"
	"github.com/stretchr/testify/assert"
)

func TestCanLoadExampleItems(t *testing.T) {
	for _, itemsFile := range []string{"items25.json", "items30.json"} {
		t.Run(itemsFile, func(t *testing.T) {
			if _, err := fs.Stat(os.DirFS("."), "examples/"+itemsFile); err != nil {
				t.Skip("no example file")
			}

			file, err := os.DirFS(".").Open("examples/" + itemsFile)
			assert.NoError(t, err)

			decoder := json.NewDecoder(file)
			decoder.DisallowUnknownFields()
			var items []api.Item
			err = decoder.Decode(&items)
			assert.NoError(t, err)

			t.Logf("loaded %d items", len(items))
		})
	}
}
