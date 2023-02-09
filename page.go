package badgrlib

import (
	"text/template"
)

var (
	G_TEMPLATE_SRC = "<g " +
		"transform=\"translate({{.Translate.X}},{{.Translate.Y}})\">" +
		"<rect width=\"{{.Dimensions.Width}}\" height=\"{{.Dimensions.Height}}\" " +
		"x=\"0\" y=\"0\" style=\"fill:none;stroke:#000000;stroke-width:0.5\">" +
		"</rect>{{.Contents}}</g>"
	G_TEMPLATE = template.Must(template.New("g").Parse(G_TEMPLATE_SRC))
)

type Translate struct {
	X float64
	Y float64
}

type GTemplateSource struct {
	Translate  Translate
	Dimensions Dimensions
	Contents   string
}

func FitObjectsOnPage(format Format, svg_objects []string) (string, error) {
	objects_on_page := format.PaperFit.X * format.PaperFit.Y
	if len(svg_objects) < objects_on_page {
		objects_on_page = len(svg_objects)
	}

	contents := ""

	for i := 0; i < objects_on_page; i++ {
		translate := Translate{
			X: format.Dimensions.Width * float64(i%format.PaperFit.X),
			Y: format.Dimensions.Height * float64(i/format.PaperFit.X),
		}

		g_data := GTemplateSource{
			Translate:  translate,
			Dimensions: format.Dimensions,
			Contents:   svg_objects[i],
		}

		g_text, err := executeTemplate(G_TEMPLATE, g_data)

		if err != nil {
			return "", err
		}

		contents += g_text
	}

	page_dimensions := Dimensions{
		Width:  format.Dimensions.Width * float64(format.PaperFit.X),
		Height: format.Dimensions.Height * float64(format.PaperFit.Y),
	}

	return WrapSvg(contents, page_dimensions)
}
