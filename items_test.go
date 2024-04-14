package main

import (
	"embed"
	"encoding/json"
	"strings"
	"testing"

	"github.com/creativeprojects/gopenhab/api"
	"github.com/stretchr/testify/assert"
)

//go:embed examples
var exampleFiles embed.FS

// load any file named item*.json in folder examples/
// and try to decode the content as []api.Item
func TestCanLoadExampleItems(t *testing.T) {
	t.Parallel()
	files, err := exampleFiles.ReadDir("examples")
	if err != nil || len(files) == 0 {
		t.Skip("no example file")
	}
	for _, itemsFile := range files {
		itemsFile := itemsFile // remove after go1.22
		if strings.HasPrefix(itemsFile.Name(), "items") && strings.HasSuffix(itemsFile.Name(), ".json") {
			t.Run(itemsFile.Name(), func(t *testing.T) {
				t.Parallel()
				file, err := exampleFiles.Open("examples/" + itemsFile.Name())
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
}
