package html

import (
	"time"

	. "github.com/maragudk/gomponents"
	. "github.com/maragudk/gomponents/html"
)

func AboutPage() Node {
	now := time.Now()

	return page("About",
		H1(Text("About")),

		P(Textf("Built with gomponents and rendered at %v.", now.Format(time.TimeOnly))),

		P(
			If(now.Second()%2 == 0, Text("It's an even second!")),
			If(now.Second()%2 != 0, Text("It's an odd second!")),
		),

		Img(Class("max-w-sm"), Src("https://www.gomponents.com/images/logo.png"), Alt("gomponents logo")),
	)
}
