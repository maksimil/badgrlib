package badgrlib

import (
	"bytes"

	"github.com/flopp/go-findfont"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/pdf"
)

type FileData = []byte

func FindArial() (*canvas.FontFamily, error) {
	arial_family := canvas.NewFontFamily("arial")
	arial_path, err := findfont.Find("arial.ttf")

	if err != nil {
		return nil, err
	}

	err = arial_family.LoadFontFile(arial_path, canvas.FontRegular)

	if err != nil {
		return nil, err
	}

	return arial_family, nil
}

func RenderPdf(
	font_family *canvas.FontFamily,
	format Format,
	page_drawers []ContextDrawer,
) (FileData, error) {
	var pdf_buffer bytes.Buffer
	page_width := format.Dimensions.Width * float64(format.PaperFit.X)
	page_height := format.Dimensions.Height * float64(format.PaperFit.Y)

	pdf_renderer := pdf.New(&pdf_buffer, page_width, page_height, &pdf.DefaultOptions)

	draw_context := Context{
		CanvasContext: canvas.NewContext(pdf_renderer), FontFamily: font_family}

	page_drawers[0](draw_context)

	for i := 1; i < len(page_drawers); i++ {
		pdf_renderer.NewPage(page_width, page_height)

		page_drawers[i](draw_context)
	}

	err := pdf_renderer.Close()

	if err != nil {
		return []byte{}, err
	}

	return pdf_buffer.Bytes(), nil
}
