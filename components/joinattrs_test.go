package components_test

import (
	"os"

	g "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// myButton is a reusable button component that accepts children and adds a "button" class.
// JoinAttrs merges all class attributes from children with the component's own class.
func myButton(children ...g.Node) g.Node {
	return Div(JoinAttrs("class", g.Group(children), Class("button")))
}

// myPrimaryButton builds on myButton, adding a "primary" class.
// The classes are merged: "primary" from here and "button" from myButton.
func myPrimaryButton(text string) g.Node {
	return myButton(Class("primary"), g.Text(text))
}

func ExampleJoinAttrs() {
	danceButton := myPrimaryButton("Dance")
	_ = danceButton.Render(os.Stdout)
	// Output: <div class="primary button">Dance</div>
}
