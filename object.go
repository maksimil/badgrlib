package badgrlib

import (
	"github.com/tdewolff/canvas"
)

type Context struct {
	CanvasContext *canvas.Context
	FontFamily    *canvas.FontFamily
}

type ContextDrawer = func(Context)

func CreateObjectDrawer(format Format, data map[string]string) ContextDrawer {
	return func(context Context) {
		for _, object := range format.Objects {
			font_face := context.FontFamily.Face(object.FontSize, canvas.Black,
				canvas.FontRegular, canvas.FontNormal)
			object_text := canvas.NewTextLine(font_face, data[object.FieldName], canvas.Left)

			context.CanvasContext.DrawText(object.X, object.Y, object_text)
		}
	}
}
