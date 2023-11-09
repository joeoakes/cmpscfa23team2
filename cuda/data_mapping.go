package cuda

// The code needs to be changed this is only an overview or a template

import (
	"github.com/kniren/gota/dataframe"
	"strings"
)

type RawData struct {
	HTMLContent string
}

func MapDataToDataSource(raw RawData) DataSource {
	return DataSource{
		Title:       "Extracted Title",
		Description: "Extracted Description",
		URL:         "https://extracted-url.com",
		Content:     "Extracted content...",
	}
}

func CleanData(raw RawData) DataSource {
	cleanedHTMLContent := strings.ReplaceAll(raw.HTMLContent, "<script>", "")
	cleanedHTMLContent = strings.ReplaceAll(cleanedHTMLContent, "</script>", "")
	df := dataframe.ReadCSV(strings.NewReader(cleanedHTMLContent))
	return DataSource{
		Title:       "Extracted Title",
		Description: "Extracted Description",
		URL:         "https://extracted-url.com",
		Content:     "Extracted content...",
	}
}
