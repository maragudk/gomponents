package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Component() Node {
	return HTML(
		Head(
			TitleEl(),
			StyleEl(),
		),
		Body(
			BlockQuote(),
			DataEl(),
			DataList(),
			FieldSet(),
			HGroup(),
			IFrame(),
			NoScript(),
			OptGroup(),
			SVG(),
			Table(
				ColGroup(),
				THead(),
				TBody(),
				TFoot(),
			),
			Figure(
				FigCaption(),
			),
		),
	)
}
