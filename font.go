package badgrlib

import (
	"math"
	"os"

	"github.com/flopp/go-findfont"
	"golang.org/x/image/font"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

var (
	FONT_FACTOR_WIDTH   = 1055.7449411764708
	FONT_FACTOR_HEIGHT  = FONT_FACTOR_WIDTH
	FONT_FACTOR_BEARING = 0.02960493827160494
)

func findArial() *sfnt.Font {
	arial_path, err := findfont.Find("arial.ttf")

	if err != nil {
		panic(err)
	}

	arial_source, err := os.ReadFile(arial_path)

	if err != nil {
		panic(err)
	}

	arial_font, err := sfnt.Parse(arial_source)

	if err != nil {
		panic(err)
	}

	return arial_font
}

func measureRune(buffer *sfnt.Buffer, used_font *sfnt.Font, r rune) (fixed.Rectangle26_6, fixed.Int26_6) {
	idx, _ := used_font.GlyphIndex(buffer, r)
	bounds, advance, _ := used_font.GlyphBounds(buffer, idx, fixed.I(1024), font.HintingNone)
	return bounds, advance
}

type StringMeasurement struct {
	Width       int
	Height      int
	LeftBearing int
}

func measureString(s string) StringMeasurement {
	arial := findArial()
	buffer := &sfnt.Buffer{}

	height := fixed.I(0)
	width := fixed.I(0)

	for _, c := range s {
		bounds, advance := measureRune(buffer, arial, c)

		glyph_height := -bounds.Min.Y
		glyph_width := advance

		if glyph_height > height {
			height = glyph_width
		}

		width = width + glyph_width
	}

	first_bounds, _ := measureRune(buffer, arial, []rune(s)[0])

	return StringMeasurement{
		Width: width.Round(), Height: height.Round(), LeftBearing: first_bounds.Min.X.Round()}
}

type StringFit struct {
	FontSize    float64
	LeftBearing float64
}

func stringFit(s string, box_width float64, box_height float64) StringFit {
	string_measurement := measureString(s)

	width_fit := box_width * FONT_FACTOR_WIDTH / float64(string_measurement.Width)

	height_fit := box_height * FONT_FACTOR_HEIGHT / float64(string_measurement.Height)

	return StringFit{
		FontSize:    math.Min(width_fit, height_fit),
		LeftBearing: float64(string_measurement.LeftBearing) * FONT_FACTOR_BEARING,
	}
}
