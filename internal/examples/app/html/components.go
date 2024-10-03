package html

import (
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

func page(title string, children ...Node) Node {
	return HTML5(HTML5Props{
		Title:    title,
		Language: "en",
		Head: []Node{
			Script(Src("https://cdn.tailwindcss.com?plugins=typography")),
		},
		Body: []Node{Class("bg-gradient-to-b from-white to-indigo-100 bg-no-repeat"),
			Div(Class("min-h-screen flex flex-col justify-between"),
				header(),
				Div(Class("grow"),
					container(true,
						Div(Class("prose prose-lg prose-indigo"),
							Group(children),
						),
					),
				),
				footer(),
			),
		},
	})
}

func header() Node {
	return Div(Class("bg-indigo-600 text-white shadow"),
		container(false,
			Div(Class("flex items-center space-x-4 h-8"),
				headerLink("/", "Home"),
				headerLink("/about", "About"),
			),
		),
	)
}

func headerLink(href, text string) Node {
	return A(Class("hover:text-indigo-300"), Href(href), Text(text))
}

func container(padY bool, children ...Node) Node {
	return Div(
		Classes{
			"max-w-7xl mx-auto":     true,
			"px-4 md:px-8 lg:px-16": true,
			"py-4 md:py-8":          padY,
		},
		Group(children),
	)
}

func footer() Node {
	return Div(Class("bg-gray-900 text-white shadow text-center h-16 flex items-center justify-center"),
		A(Href("https://www.gomponents.com"), Text("gomponents")),
	)
}
