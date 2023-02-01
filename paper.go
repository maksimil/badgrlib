package badgrlib

import (
	"text/template"
)

var (
	SVG_TEMPLATE_SRC = "<svg width=\"{{.Dimensions.Width}}mm\" height=\"{{.Dimensions.Height}}mm\" " +
		"viewBox=\"0 0 {{.Dimensions.Width}} {{.Dimensions.Height}}\">{{.Contents}}</svg>"
	SVG_TEMPLATE = template.Must(template.New("svg").Parse(SVG_TEMPLATE_SRC))

	G_TEMPLATE_SRC = "<g transform=\"translate({{.Translate.X}}," +
		"{{.Translate.Y}})\">{{.Contents}}</g>"
	G_TEMPLATE = template.Must(template.New("g").Parse(G_TEMPLATE_SRC))
)

type SvgTemplateSource struct {
	Dimensions Dimensions
	Contents   string
}

type Translate struct {
	X float64
	Y float64
}

type GTemplateSource struct {
	Translate Translate
	Contents  string
}

func FitSvgsToPaper(format Format, svg_strings []string) (string, error) {
	objects_on_page := format.PaperFit.X * format.PaperFit.Y
	if len(svg_strings) < objects_on_page {
		objects_on_page = len(svg_strings)
	}

	contents := ""

	for i := 0; i < objects_on_page; i++ {
		x_pos := float64(i % format.PaperFit.X)
		y_pos := float64(i / format.PaperFit.X)

		group_data := GTemplateSource{
			Translate: Translate{X: x_pos * format.Dimensions.Width, Y: y_pos * format.Dimensions.Height},
			Contents:  svg_strings[i],
		}

		group_string, err := executeTemplate(G_TEMPLATE, group_data)

		if err != nil {
			return "", err
		}

		contents += group_string
	}

	svg_data := SvgTemplateSource{
		Dimensions: Dimensions{
			Width:  format.Dimensions.Width * float64(format.PaperFit.X),
			Height: format.Dimensions.Height * float64(format.PaperFit.Y),
		},
		Contents: contents,
	}

	svg_string, err := executeTemplate(SVG_TEMPLATE, svg_data)

	if err != nil {
		return "", err
	}

	return svg_string, nil
}
