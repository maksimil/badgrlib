package badgrlib

import (
	"bytes"
	"text/template"
)

var (
	TEXT_TEMPLATE_SRC = "<text x=\"{{.X}}\" y=\"{{.Y}}\" " +
		"style=\"font-size:{{.FontSize}}px;line-height:1.25;font-family:Arial\">" +
		"{{.Text}}</text>"
	TEXT_TEMPLATE = template.Must(template.New("text").Parse(TEXT_TEMPLATE_SRC))

	SVG_TEMPLATE_SRC = "<svg " +
		"width=\"{{.Dimensions.Width}}mm\" height=\"{{.Dimensions.Height}}mm\" " +
		"viewBox=\"0 0 {{.Dimensions.Width}} {{.Dimensions.Height}}\">{{.Contents}}</svg>"
	SVG_TEMPLATE = template.Must(template.New("svg").Parse(SVG_TEMPLATE_SRC))
)

type TextTemplateSource struct {
	X        float64
	Y        float64
	FontSize float64
	Text     string
}

type SvgTemplateSource struct {
	Dimensions Dimensions
	Contents   string
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, data)
	return buf.String(), err
}

func CreateSingleSvg(format Format, data map[string]string) (string, error) {
	contents := ""

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

		contents += text_contents
	}

	return contents, nil
}

func WrapSvg(svg string, dimensions Dimensions) (string, error) {
	svg_data := SvgTemplateSource{
		Dimensions: dimensions,
		Contents:   svg,
	}
	return executeTemplate(SVG_TEMPLATE, svg_data)
}
