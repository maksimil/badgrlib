package badgrlib

import (
	"bytes"
	"text/template"
)

var (
	BOX_TEMPLATE_SRC = "<rect width=\"{{.Dimensions.Width}}\" height=\"{{.Dimensions.Height}}\" " +
		"x=\"{{.X}}\" y=\"{{.Y}}\" style=\"fill:none;stroke:#000000;stroke-width:0.5\"/>"
	BOX_TEMPLATE = template.Must(template.New("box").Parse(BOX_TEMPLATE_SRC))

	TEXT_TEMPLATE_SRC = "<text x=\"{{.X}}\" y=\"{{.Y}}\" " +
		"style=\"font-size:{{.FontSize}}px;line-height:1.25;font-family:Arial\">" +
		"{{.Text}}</text>"
	TEXT_TEMPLATE = template.Must(template.New("text").Parse(TEXT_TEMPLATE_SRC))
)

type BoxTemplateSource struct {
	X          float64
	Y          float64
	Dimensions Dimensions
}

type TextTemplateSource struct {
	X        float64
	Y        float64
	FontSize float64
	Text     string
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, data)
	return buf.String(), err
}

func CreateSingleSvg(format Format, data map[string]string) (string, error) {
	contents := ""

	box_data := BoxTemplateSource{
		X:          0,
		Y:          0,
		Dimensions: format.Dimensions,
	}

	box_contents, err := executeTemplate(BOX_TEMPLATE, box_data)

	if err != nil {
		return "", err
	}

	for _, object := range format.Objects {
		text_data := TextTemplateSource{
			X:        object.X,
			Y:        object.Y,
			FontSize: object.FontSize,
			Text:     data[object.FieldName],
		}

		text_contents, err := executeTemplate(TEXT_TEMPLATE, text_data)

		if err != nil {
			return "", err
		}

		contents += box_contents + text_contents
	}

	return contents, nil
}
