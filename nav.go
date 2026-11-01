package bs

import (
	"strconv"

	"github.com/deltegui/hx"
)

// Tab is an entry of Tabs.
type Tab struct{ ID, Label string }

// Tabs is a row of tabs. active holds the ID of the chosen one and clicking a tab
// sets it; show the content with TabPanel. Mods: Pills makes them pills, Underline
// a line.
func Tabs(active *hx.SignalT[string], tabs []Tab, mods ...Mod) hx.INode {
	kind := "nav-tabs"
	for _, m := range mods {
		kind = "nav-" + string(m)
	}
	items := make([]hx.INode, len(tabs))
	for i, t := range tabs {
		b := hx.Button().Type("button").Class("nav-link").Attribute("role", "tab").Text(t.Label)
		b.OnClick(func(hx.EventContext) { active.Set(t.ID) })
		b.Effect(func() {
			on := active.Get() == t.ID
			b.ClassIf("active", on)
			b.Attribute("aria-selected", strconv.FormatBool(on))
		})
		items[i] = hx.Li().Class("nav-item").Attribute("role", "presentation").Body(b)
	}
	return hx.Ul().Class("nav", kind).Attribute("role", "tablist").BodyList(items)
}

// Tabs variants.
const (
	Pills     Mod = "pills"
	Underline Mod = "underline"
)

// TabPanel shows build() while the Tabs have id chosen.
func TabPanel(active hx.Gettable[string], id string, build func() hx.INode) hx.INode {
	return hx.Dynamic(func() hx.INode {
		if active.Get() != id {
			return nil
		}
		var panel hx.INode
		hx.Untrack(func() { panel = build() })
		return hx.Div().Attribute("role", "tabpanel").Body(panel)
	})
}

// Nav is a list of links or buttons (NavLink). Chain "flex-column" for a
// vertical one, "justify-content-center" to center it.
func Nav(children ...hx.INode) hx.INode { return hx.Ul().Class("nav").Body(children...) }

// NavLink is an item of a Nav. Chain "active" or "disabled" on the returned
// node; the class goes to the link, not to the item.
func NavLink(label, href string, active bool) hx.INode {
	a := hx.A().Class("nav-link").Href(href).Text(label)
	if active {
		a.Class("active").Attribute("aria-current", "page")
	}
	return hx.Li().Class("nav-item").Body(a)
}

// Navbar is the bar at the top of the page. expand is the breakpoint from which
// it is not folded ("lg"); put a NavbarBrand, a NavbarToggler and a
// NavbarCollapse in it, inside a Container. Chain "bg-body-tertiary" or
// "bg-dark" and "navbar-dark" for colors.
func Navbar(expand string, children ...hx.INode) hx.INode {
	return hx.Nav().Class("navbar", "navbar-expand-"+expand).Body(children...)
}

// NavbarBrand is the logo or the name of the site.
func NavbarBrand(label, href string) hx.INode {
	return hx.A().Class("navbar-brand").Href(href).Text(label)
}

// NavbarToggler is the button that opens the menu on a phone. It toggles open,
// which NavbarCollapse reads.
func NavbarToggler(open *hx.SignalT[bool]) hx.INode {
	b := hx.Button().Type("button").Class("navbar-toggler").Attribute("aria-label", "Toggle navigation").
		Body(hx.Span().Class("navbar-toggler-icon"))
	b.OnClick(func(hx.EventContext) { open.Set(!open.Get()) })
	b.Effect(func() { b.Attribute("aria-expanded", strconv.FormatBool(open.Get())) })
	return b
}

// NavbarCollapse is the part that folds on a phone.
func NavbarCollapse(open hx.Gettable[bool], children ...hx.INode) hx.INode {
	return BindClass(hx.Div().Class("collapse", "navbar-collapse").Body(children...), "show", open)
}

// NavbarNav is the list of links of a NavbarCollapse. Chain "me-auto" to push
// what follows to the right.
func NavbarNav(children ...hx.INode) hx.INode { return hx.Ul().Class("navbar-nav").Body(children...) }

// NavbarLink is a link of the navbar.
func NavbarLink(label, href string, active bool) hx.INode {
	return NavLink(label, href, active)
}

// Carousel shows one of its slides at a time. index is the slide that shows
// (from 0) and the arrows and the dots change it; it wraps around. It does not
// slide: Bootstrap does that with JavaScript. Each slide is a CarouselItem.
func Carousel(index *hx.SignalT[int], items ...hx.INode) hx.INode {
	n := len(items)
	for i, it := range items {
		it.Effect(func() { it.ClassIf("active", index.Get() == i) })
	}
	dots := make([]hx.INode, n)
	for i := range dots {
		d := hx.Button().Type("button").Attribute("aria-label", "Slide "+strconv.Itoa(i+1))
		d.OnClick(func(hx.EventContext) { index.Set(i) })
		d.Effect(func() { d.ClassIf("active", index.Get() == i) })
		dots[i] = d
	}
	step := func(class, icon, label string, delta int) hx.INode {
		return hx.Button().Type("button").Class(class).OnClick(func(hx.EventContext) {
			if n > 0 {
				index.Set((index.Get() + delta + n) % n)
			}
		}).Body(
			hx.Span().Class(icon).Attribute("aria-hidden", "true"),
			hx.Span().Class("visually-hidden").Text(label))
	}
	return hx.Div().Class("carousel", "slide").Body(
		hx.Div().Class("carousel-indicators").BodyList(dots),
		hx.Div().Class("carousel-inner").BodyList(items),
		step("carousel-control-prev", "carousel-control-prev-icon", "Previous", -1),
		step("carousel-control-next", "carousel-control-next-icon", "Next", 1))
}

// CarouselItem is a slide: usually an image, with a caption if you want one.
func CarouselItem(children ...hx.INode) hx.INode {
	return hx.Div().Class("carousel-item").Body(children...)
}

// CarouselCaption is the text over a slide.
func CarouselCaption(title, text string) hx.INode {
	return hx.Div().Class("carousel-caption", "d-none", "d-md-block").Body(
		hx.H5().Text(title), hx.P().Text(text))
}
