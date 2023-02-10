package badgrlib

import (
	"github.com/tdewolff/canvas"
)

func FitObjectsOnPage(format Format, drawers []ContextDrawer) ContextDrawer {
	return func(context Context) {
		context.CanvasContext.SetCoordSystem(canvas.CartesianIV)
		context.CanvasContext.DrawPath(0, 0,
			canvas.Grid(
				format.Dimensions.Width*float64(format.PaperFit.X),
				format.Dimensions.Height*float64(format.PaperFit.Y),
				format.PaperFit.X, format.PaperFit.Y, 0.5))

		objects_on_page := format.PaperFit.X * format.PaperFit.Y

		if len(drawers) < objects_on_page {
			objects_on_page = len(drawers)
		}

		for i := 0; i < objects_on_page; i++ {
			context.CanvasContext.
				Translate(
					format.Dimensions.Width*float64(i%format.PaperFit.X),
					-format.Dimensions.Height*float64(i/format.PaperFit.X),
				)

			drawers[i](context)

			context.CanvasContext.ResetView()
		}
	}
}
