package bs_test

import (
	"strings"
	"testing"

	"github.com/deltegui/hx"
	"github.com/deltegui/hx/ssr"
	"github.com/deltegui/hxbs"
)

func render(build func() hx.INode) string { return ssr.Render(build) }

func TestComponents(t *testing.T) {
	text := hx.Signal("hi")
	flag := hx.Signal(true)
	off := hx.Signal(false)
	num := hx.Signal(3)
	tab := hx.Signal("b")
	acc := hx.Signal("one")
	choice := hx.Signal("x")
	nop := func() {}
	node := func() hx.INode { return hx.P().Text("x") }

	cases := []struct {
		name  string
		build func() hx.INode
		want  []string
	}{
		{"button", func() hx.INode { return bs.Button("Go", bs.Danger, bs.Lg) }, []string{`class="btn btn-danger btn-lg"`, `type="button"`, `>Go</button>`}},
		{"button default and outline", func() hx.INode {
			return hx.Div().Body(bs.Button("a"), bs.Button("b", bs.Outline(bs.Success)), bs.SubmitButton("c"), bs.ButtonLink("/x", "d", bs.Link))
		}, []string{`btn btn-primary`, `btn btn-outline-success`, `type="submit"`, `btn btn-link`, `href="/x"`}},
		{"button group", func() hx.INode { return bs.ButtonGroup([]bs.Mod{bs.Vertical}, bs.Button("a")) }, []string{`btn-group-vertical`}},
		{"alert", func() hx.INode { return bs.DismissibleAlert(nop, node(), bs.Warning) }, []string{`alert alert-dismissible alert-warning`, `role="alert"`, `btn-close`}},
		{"badge", func() hx.INode { return bs.Badge("new", bs.Info, bs.Pill) }, []string{`badge`, `text-bg-info`, `rounded-pill`}},
		{"breadcrumb", func() hx.INode {
			return bs.Breadcrumb(bs.Crumb{Label: "Home", Href: "/"}, bs.Crumb{Label: "Now", Current: true})
		}, []string{`<ol class="breadcrumb">`, `<li class="breadcrumb-item"><a href="/">Home`, `active breadcrumb-item`, `aria-current="page"`}},
		{"card", func() hx.INode {
			return bs.Card(bs.CardImage("/a", "a"), bs.CardBody(bs.CardTitle("T"), bs.CardText("x")), bs.CardFooter(node()))
		}, []string{`class="card"`, `card-img-top`, `card-body`, `card-title">T`, `card-text`, `card-footer`}},
		{"list group", func() hx.INode {
			return bs.ListGroup([]bs.Mod{bs.Flush}, bs.ListGroupItem(node(), bs.Success), bs.ListGroupAction("go", nop))
		}, []string{`list-group list-group-flush`, `list-group-item list-group-item-success`, `list-group-item-action`}},
		{"pagination", func() hx.INode { return bs.Pagination(bs.PaginationOptions{Page: hx.Signal(5), Pages: hx.Signal(10)}) }, []string{`class="pagination"`, `active page-item`, `aria-current="page"`, `…`}},
		{"pagination first", func() hx.INode {
			return bs.Pagination(bs.PaginationOptions{Page: hx.Signal(1), Pages: hx.Signal(1), OnPage: func(int) {}, Mods: []bs.Mod{bs.Sm}})
		}, []string{`pagination-sm`, `disabled page-item`}},
		{"progress", func() hx.INode { return bs.ProgressOf(num, 10, bs.Success, bs.Animated) }, []string{`role="progressbar"`, `aria-valuenow="3"`, `width:30%`, `progress-bar-animated`, `text-bg-success`}},
		{"spinner", func() hx.INode { return hx.Div().Body(bs.Spinner(bs.Primary, bs.Sm), bs.Spinner(bs.Grow)) }, []string{`spinner-border spinner-border-sm text-primary`, `spinner-grow`, `visually-hidden`}},
		{"placeholder", func() hx.INode { return bs.PlaceholderGlow(bs.Placeholder(6)) }, []string{`placeholder-glow`, `col-6 placeholder`}},
		{"layout", func() hx.INode {
			return bs.Container(bs.Row(bs.ColOn(bs.MD, 6, node()).Class(bs.Offset(2)), bs.Col(), bs.ColSpan(3)))
		}, []string{`class="container"`, `class="row"`, `col-md-6 offset-2`, `class="col"`, `col-3`}},
		{"stack ratio", func() hx.INode { return bs.HStack(bs.Ratio("16x9", node())).Class(bs.Gap(3)) }, []string{`gap-3 hstack`, `ratio ratio-16x9`}},
		{"typography", func() hx.INode { return hx.Div().Body(bs.Display(2, "T"), bs.Lead("L"), bs.Heading(3, "H")) }, []string{`display-2`, `class="lead"`, `<h3>H</h3>`}},
		{"figure", func() hx.INode { return bs.Figure(bs.Image("/a", "A"), "cap") }, []string{`figure-img img-fluid rounded`, `figure-caption`}},
		{"table", func() hx.INode {
			return bs.TableResponsive(bs.SimpleTable([]bs.Mod{bs.Striped, bs.Hover}, []string{"A", "B"}, hx.Tr()))
		}, []string{`table-responsive`, `table table-hover table-striped`, `<th scope="col">A</th>`}},
		{"form field", func() hx.INode {
			return bs.FormField("Name", "name", bs.Input("text", "name", text, bs.Lg), hx.Signal([]string{"bad", "worse"}))
		}, []string{`class="mb-3"`, `for="name"`, `form-label`, `form-control form-control-lg`, `is-invalid`, `invalid-feedback`, `bad · worse`}},
		{"form controls", func() hx.INode {
			return hx.Div().Body(bs.TextArea("t", text), bs.Select("s", text, "—", bs.Plain("a", "b")), bs.Checkbox("c", "ok", off), bs.Switch("sw", "on", off),
				bs.Radios("r", choice, bs.Plain("x", "y")), bs.FileInput("f", nil), bs.Range("rg", text))
		}, []string{`<textarea`, `class="form-select"`, `<option>—</option>`, `form-check-input`, `form-check form-switch`, `type="radio"`, `for="r-x"`, `form-control`, `form-range`}},
		{"floating", func() hx.INode { return bs.FloatingLabel("Email", "e", bs.Input("email", "e", text)) }, []string{`form-floating`, `placeholder="Email"`, `<label for="e">Email`}},
		{"input group", func() hx.INode {
			return bs.InputGroup(bs.InputGroupText(hx.Span().Text("@")), bs.Input("text", "u", text))
		}, []string{`input-group`, `input-group-text`}},
		{"modal open", func() hx.INode {
			return bs.Modal(bs.ModalOptions{Open: flag, Title: "Title", OnClose: nop, Body: node, Footer: node, Size: bs.Lg})
		}, []string{`fade modal-backdrop show`, `d-block fade modal show`, `modal-dialog modal-lg`, `modal-title">Title`, `modal-body`, `modal-footer`, `btn-close`}},
		{"modal centered", func() hx.INode { return bs.Modal(bs.ModalOptions{Open: flag, Title: "T", Body: node, Centered: true}) }, []string{`modal-dialog-centered`}},
		{"modal closed", func() hx.INode { return hx.Div().Body(bs.Modal(bs.ModalOptions{Open: off, Title: "T", Body: node})) }, []string{`<div></div>`}},
		{"offcanvas", func() hx.INode {
			return bs.Offcanvas(bs.OffcanvasOptions{Open: flag, Side: "end", Title: "Menu", Body: node})
		}, []string{`offcanvas-backdrop`, `offcanvas offcanvas-end show`, `offcanvas-title">Menu`, `offcanvas-body`}},
		{"dropdown open", func() hx.INode {
			return bs.Dropdown(flag, "Menu", []hx.INode{bs.DropdownHeader("H"), bs.DropdownItem("One", "/1", nil), bs.DropdownDivider()}, bs.Info)
		}, []string{`class="dropdown"`, `btn btn-info dropdown-toggle`, `aria-expanded="true"`, `dropdown-menu show`, `dropdown-header`, `dropdown-item`, `dropdown-divider`, `position:fixed`}},
		{"dropdown closed", func() hx.INode { return bs.Dropdown(off, "Menu", nil) }, []string{`btn-secondary`, `aria-expanded="false"`, `class="dropdown-menu"`}},
		{"collapse", func() hx.INode { return hx.Div().Body(bs.Collapse(flag, node()), bs.Collapse(off, node())) }, []string{`collapse show`, `class="collapse"`}},
		{"accordion", func() hx.INode {
			return bs.Accordion(acc, bs.AccordionItem{ID: "one", Title: "One", Body: node}, bs.AccordionItem{ID: "two", Title: "Two", Body: node})
		}, []string{`accordion-item`, `accordion-button"`, `accordion-button collapsed`, `accordion-collapse collapse show`, `accordion-body`}},
		{"toast", func() hx.INode {
			return bs.ToastContainer("top-0 end-0", bs.Toast(bs.ToastOptions{Open: flag, Title: "Saved", Body: node}))
		}, []string{`end-0 p-3 position-fixed toast-container top-0`, `show toast`, `toast-header`, `toast-body`}},
		{"tabs", func() hx.INode {
			return hx.Div().Body(bs.Tabs(tab, []bs.Tab{{"a", "A"}, {"b", "B"}}, bs.Pills),
				bs.TabPanel(tab, "a", func() hx.INode { return hx.P().Text("PANEL-A") }),
				bs.TabPanel(tab, "b", func() hx.INode { return hx.P().Text("PANEL-B") }))
		}, []string{`nav nav-pills`, `active nav-link`, `aria-selected="true"`, `PANEL-B`}},
		{"nav", func() hx.INode { return bs.Nav(bs.NavLink("A", "/a", true)) }, []string{`class="nav"`, `active nav-link`, `aria-current="page"`}},
		{"navbar", func() hx.INode {
			return bs.Navbar("lg", bs.Container(bs.NavbarBrand("Logo", "/"), bs.NavbarToggler(flag),
				bs.NavbarCollapse(flag, bs.NavbarNav(bs.NavbarLink("A", "/a", false)).Class("me-auto"))))
		}, []string{`navbar navbar-expand-lg`, `navbar-brand`, `navbar-toggler`, `aria-expanded="true"`, `collapse navbar-collapse show`, `navbar-nav`}},
		{"carousel", func() hx.INode {
			return bs.Carousel(hx.Signal(1), bs.CarouselItem(node()), bs.CarouselItem(node(), bs.CarouselCaption("T", "x")))
		}, []string{`carousel slide`, `carousel-indicators`, `active carousel-item`, `carousel-control-prev`, `carousel-control-next-icon`}},
		{"helpers", func() hx.INode {
			return hx.Div().Class(bs.TextColor(bs.Primary), bs.BgColor(bs.Light), bs.Space("m", "t", 3), bs.TextBg(bs.Dark), bs.BorderColor(bs.Danger))
		}, []string{`text-primary`, `bg-light`, `mt-3`, `text-bg-dark`, `border-danger`}},
		{"icon", func() hx.INode { return bs.Icon("house") }, []string{`bi bi-house`}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			html := render(c.build)
			for _, w := range c.want {
				if !strings.Contains(html, w) {
					t.Errorf("missing %q in\n%s", w, html)
				}
			}
		})
	}
	if html := render(func() hx.INode { return hx.Div().Body(bs.TabPanel(tab, "a", node)) }); strings.Contains(html, "<p>") {
		t.Errorf("TabPanel shows an inactive tab: %s", html)
	}
	if html := render(func() hx.INode {
		return bs.Accordion(acc, bs.AccordionItem{ID: "one", Title: "One"}, bs.AccordionItem{ID: "two", Title: "Two"})
	}); strings.Count(html, "collapse show") != 1 {
		t.Errorf("accordion must open one section: %s", html)
	}
}
