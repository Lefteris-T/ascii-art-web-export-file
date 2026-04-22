package export

import (
	"encoding/json"
	"fmt"
	"html/template"
)

type Payload struct {
	Text   string `json:"text"`
	Banner string `json:"banner"`
	Result string `json:"result"`
}

func Build(format, text, banner, result string) (content string, contentType string, filename string, err error) {
	switch format {
	case "txt":
		return result, "text/plain; charset=utf-8", "ascii-art.txt", nil

	case "html":
		html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>ASCII Art Export</title>
</head>
<body>
  <pre>%s</pre>
</body>
</html>`, template.HTMLEscapeString(result))

		return html, "text/html; charset=utf-8", "ascii-art.html", nil

	case "json":
		payload := Payload{
			Text:   text,
			Banner: banner,
			Result: result,
		}

		data, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return "", "", "", err
		}

		return string(data), "application/json", "ascii-art.json", nil

	default:
		return "", "", "", fmt.Errorf("unsupported export format")
	}
}
