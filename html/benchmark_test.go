//go:build go1.24

package html_test

import (
	"io"
	"testing"

	g "maragu.dev/gomponents"
	c "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
)

// BenchmarkRealisticPage benchmarks rendering a full, realistic HTML page
// resembling a typical web application dashboard with navigation, sidebar,
// content cards, a data table, and a footer.
func BenchmarkRealisticPage(b *testing.B) {
	type navItem struct {
		href, label string
		active      bool
	}

	type card struct {
		title, description, badge string
		count                     int
	}

	type row struct {
		name, email, role, status string
	}

	navItems := []navItem{
		{"/dashboard", "Dashboard", true},
		{"/projects", "Projects", false},
		{"/team", "Team", false},
		{"/reports", "Reports", false},
		{"/settings", "Settings", false},
	}

	cards := []card{
		{"Total Revenue", "Up 12% from last month", "success", 48250},
		{"Active Users", "Up 8% from last week", "info", 2340},
		{"New Orders", "Down 3% from yesterday", "warning", 156},
		{"Conversion Rate", "Stable", "neutral", 3},
	}

	rows := make([]row, 50)
	for i := range rows {
		rows[i] = row{
			name:   "User Name",
			email:  "user@example.com",
			role:   "Member",
			status: "Active",
		}
	}

	isAdmin := true

	page := func() g.Node {
		return c.HTML5(c.HTML5Props{
			Title:       "Dashboard - My Application",
			Description: "Application dashboard with analytics and user management.",
			Language:    "en",
			Head: g.Group{
				Link(Rel("stylesheet"), Href("/css/app.css")),
				Link(Rel("icon"), Type("image/svg+xml"), Href("/favicon.svg")),
				Meta(g.Attr("name", "theme-color"), Content("#4f46e5")),
				Script(Src("/js/app.js"), Defer()),
			},
			Body: g.Group{
				Class("min-h-screen bg-gray-50 text-gray-900 antialiased"),

				// Navigation bar
				Nav(Class("bg-indigo-600 text-white shadow-lg"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8"),
						Div(Class("flex items-center justify-between h-16"),
							// Logo
							Div(Class("flex items-center"),
								A(Href("/"), Class("text-xl font-bold tracking-tight"),
									g.Text("MyApp"),
								),
							),
							// Nav links
							Div(Class("hidden md:flex items-center space-x-4"),
								g.Map(navItems, func(item navItem) g.Node {
									return A(
										Href(item.href),
										c.Classes{
											"px-3 py-2 rounded-md text-sm font-medium": true,
											"bg-indigo-700 text-white":                  item.active,
											"text-indigo-100 hover:bg-indigo-500":        !item.active,
										},
										g.Text(item.label),
									)
								}),
							),
							// User menu
							Div(Class("flex items-center space-x-3"),
								Span(Class("text-sm"), g.Text("admin@example.com")),
								Button(Class("p-1 rounded-full hover:bg-indigo-500"),
									Type("button"),
									Span(Class("sr-only"), g.Text("Notifications")),
									g.Raw(`<svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/></svg>`),
								),
								A(Href("/logout"), Class("text-sm underline"), g.Text("Log out")),
							),
						),
					),
				),

				// Main layout: sidebar + content
				Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8"),
					Div(Class("flex gap-8"),

						// Sidebar
						Aside(Class("hidden lg:block w-64 flex-shrink-0"),
							Div(Class("bg-white rounded-lg shadow p-4"),
								H3(Class("text-sm font-semibold text-gray-500 uppercase tracking-wider mb-3"),
									g.Text("Quick Links"),
								),
								Ul(Class("space-y-2"),
									Li(A(Href("#"), Class("block px-3 py-2 rounded hover:bg-gray-100 text-sm"), g.Text("Overview"))),
									Li(A(Href("#"), Class("block px-3 py-2 rounded hover:bg-gray-100 text-sm"), g.Text("Analytics"))),
									Li(A(Href("#"), Class("block px-3 py-2 rounded hover:bg-gray-100 text-sm"), g.Text("Exports"))),
									Li(A(Href("#"), Class("block px-3 py-2 rounded hover:bg-gray-100 text-sm"), g.Text("Integrations"))),
									g.If(isAdmin,
										Li(A(Href("/admin"), Class("block px-3 py-2 rounded hover:bg-gray-100 text-sm text-indigo-600 font-medium"), g.Text("Admin Panel"))),
									),
								),
							),
						),

						// Main content
						Main(Class("flex-1 min-w-0"),
							// Page header
							Div(Class("mb-6"),
								H1(Class("text-2xl font-bold text-gray-900"), g.Text("Dashboard")),
								P(Class("mt-1 text-sm text-gray-500"), g.Text("Welcome back. Here's what's happening with your projects.")),
							),

							// Stats cards
							Div(Class("grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8"),
								g.Map(cards, func(cd card) g.Node {
									return Div(Class("bg-white rounded-lg shadow p-5"),
										Div(Class("flex items-center justify-between"),
											Div(
												P(Class("text-sm font-medium text-gray-500"), g.Text(cd.title)),
												P(Class("mt-1 text-3xl font-semibold text-gray-900"), g.Textf("%d", cd.count)),
											),
											Span(
												c.Classes{
													"inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium": true,
													"bg-green-100 text-green-800": cd.badge == "success",
													"bg-blue-100 text-blue-800":   cd.badge == "info",
													"bg-yellow-100 text-yellow-800": cd.badge == "warning",
													"bg-gray-100 text-gray-800":     cd.badge == "neutral",
												},
												g.Text(cd.description),
											),
										),
									)
								}),
							),

							// Data table
							Div(Class("bg-white shadow rounded-lg overflow-hidden"),
								Div(Class("px-4 py-5 sm:px-6 border-b border-gray-200"),
									Div(Class("flex items-center justify-between"),
										H2(Class("text-lg font-medium text-gray-900"), g.Text("Team Members")),
										Div(Class("flex space-x-2"),
											Input(Type("search"), Class("border rounded-md px-3 py-1.5 text-sm"), g.Attr("placeholder", "Search...")),
											Button(Type("button"), Class("bg-indigo-600 text-white px-4 py-1.5 rounded-md text-sm font-medium hover:bg-indigo-700"),
												g.Text("Add Member"),
											),
										),
									),
								),
								Table(Class("min-w-full divide-y divide-gray-200"),
									THead(
										Tr(Class("bg-gray-50"),
											Th(Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), g.Attr("scope", "col"), g.Text("Name")),
											Th(Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), g.Attr("scope", "col"), g.Text("Email")),
											Th(Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), g.Attr("scope", "col"), g.Text("Role")),
											Th(Class("px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"), g.Attr("scope", "col"), g.Text("Status")),
											g.If(isAdmin,
												Th(Class("px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider"), g.Attr("scope", "col"), g.Text("Actions")),
											),
										),
									),
									TBody(Class("bg-white divide-y divide-gray-200"),
										g.Map(rows, func(r row) g.Node {
											return Tr(Class("hover:bg-gray-50"),
												Td(Class("px-6 py-4 whitespace-nowrap"),
													Div(Class("flex items-center"),
														Div(Class("h-8 w-8 rounded-full bg-indigo-100 flex items-center justify-center"),
															Span(Class("text-sm font-medium text-indigo-600"), g.Text("U")),
														),
														Div(Class("ml-4"),
															Div(Class("text-sm font-medium text-gray-900"), g.Text(r.name)),
														),
													),
												),
												Td(Class("px-6 py-4 whitespace-nowrap text-sm text-gray-500"), g.Text(r.email)),
												Td(Class("px-6 py-4 whitespace-nowrap text-sm text-gray-500"), g.Text(r.role)),
												Td(Class("px-6 py-4 whitespace-nowrap"),
													Span(Class("inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"),
														g.Text(r.status),
													),
												),
												g.If(isAdmin,
													Td(Class("px-6 py-4 whitespace-nowrap text-right text-sm"),
														A(Href("#"), Class("text-indigo-600 hover:text-indigo-900 mr-3"), g.Text("Edit")),
														A(Href("#"), Class("text-red-600 hover:text-red-900"), g.Text("Delete")),
													),
												),
											)
										}),
									),
								),
								// Pagination
								Div(Class("px-4 py-3 border-t border-gray-200 sm:px-6"),
									Div(Class("flex items-center justify-between"),
										P(Class("text-sm text-gray-700"),
											g.Text("Showing 1 to 50 of 237 results"),
										),
										Div(Class("flex space-x-1"),
											Button(Type("button"), Class("px-3 py-1 border rounded text-sm"), Disabled(), g.Text("Previous")),
											Button(Type("button"), Class("px-3 py-1 border rounded text-sm bg-indigo-600 text-white"), g.Text("1")),
											Button(Type("button"), Class("px-3 py-1 border rounded text-sm"), g.Text("2")),
											Button(Type("button"), Class("px-3 py-1 border rounded text-sm"), g.Text("3")),
											Button(Type("button"), Class("px-3 py-1 border rounded text-sm"), g.Text("Next")),
										),
									),
								),
							),
						),
					),
				),

				// Footer
				Footer(Class("bg-gray-900 text-gray-400 mt-12"),
					Div(Class("max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8"),
						Div(Class("grid grid-cols-1 md:grid-cols-3 gap-8"),
							Div(
								H3(Class("text-white font-semibold mb-3"), g.Text("MyApp")),
								P(Class("text-sm"), g.Text("Building the future of project management, one dashboard at a time.")),
							),
							Div(
								H3(Class("text-white font-semibold mb-3"), g.Text("Links")),
								Ul(Class("space-y-1 text-sm"),
									Li(A(Href("#"), Class("hover:text-white"), g.Text("Documentation"))),
									Li(A(Href("#"), Class("hover:text-white"), g.Text("API Reference"))),
									Li(A(Href("#"), Class("hover:text-white"), g.Text("Status Page"))),
								),
							),
							Div(
								H3(Class("text-white font-semibold mb-3"), g.Text("Legal")),
								Ul(Class("space-y-1 text-sm"),
									Li(A(Href("#"), Class("hover:text-white"), g.Text("Privacy Policy"))),
									Li(A(Href("#"), Class("hover:text-white"), g.Text("Terms of Service"))),
								),
							),
						),
						Div(Class("mt-8 pt-8 border-t border-gray-700 text-sm text-center"),
							P(g.Textf("© %d MyApp. All rights reserved.", 2026)),
						),
					),
				),
			},
		})
	}

	b.Run("construct and render", func(b *testing.B) {
		for b.Loop() {
			_ = page().Render(io.Discard)
		}
	})

	b.Run("render only", func(b *testing.B) {
		p := page()
		for b.Loop() {
			_ = p.Render(io.Discard)
		}
	})
}
