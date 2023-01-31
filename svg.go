package badgrlib

import (
	"bytes"
	"text/template"
)

var (
	SVG_TEMPLATE_SRC = "<svg width=\"{{.Width}}mm\" height=\"{{.Height}}mm\" " +
		"viewBox=\"0 0 {{.Width}} {{.Height}}\">{{.Contents}}</svg>"
	SVG_TEMPLATE = template.Must(template.New("svg").Parse(SVG_TEMPLATE_SRC))

	BOX_TEMPLATE_SRC = "<rect width=\"{{.Width}}\" height=\"{{.Height}}\" " +
		"x=\"{{.X}}\" y=\"{{.Y}}\" style=\"fill:none;stroke:#000000;stroke-width:0.3\"/>"
	BOX_TEMPLATE = template.Must(template.New("box").Parse(BOX_TEMPLATE_SRC))
)

type BoxTemplateSource struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type SvgTemplateSource struct {
	Width    float32
	Height   float32
	Contents string
}

func executeTemplate(tmpl *template.Template, data interface{}) (string, error) {
	buf := bytes.Buffer{}
	err := tmpl.Execute(&buf, data)
	return buf.String(), err
}

func CreateSingleSvg(format Format, data map[string]string) (string, error) {
	contents := ""

	for _, object := range format.Objects {
		box_data := BoxTemplateSource{
			X:      object.X,
			Y:      object.Y,
			Width:  object.Width,
			Height: object.Height,
		}

		box_contents, err := executeTemplate(BOX_TEMPLATE, box_data)

		if err != nil {
			return "", err
		}

		contents += box_contents
	}

	svg_data := SvgTemplateSource{
		Width:    format.Dimensions.Width,
		Height:   format.Dimensions.Height,
		Contents: contents,
	}

	svg_contents, err := executeTemplate(SVG_TEMPLATE, svg_data)

	if err != nil {
		return "", err
	}

	return svg_contents, nil
}
