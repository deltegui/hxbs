package bs_test

import (
	"strings"
	"testing"

	"github.com/deltegui/hx"
	"github.com/deltegui/hxbs"
)

func TestMoreComponents(t *testing.T) {
	flag := hx.Signal(true)
	off := hx.Signal(false)
	text := hx.Signal("a")
	tab := hx.Signal("x")
	n := hx.Signal(20)
	node := func() hx.INode { return hx.P().Text("x") }
	items := []hx.INode{bs.DropdownItem("One", "", nil)}

	cases := []struct {
		name  string
		build func() hx.INode
		want  []string
		not   []string
	}{
		{"flex", func() hx.INode {
			return bs.Flex(node(), bs.Spacer(), node()).Class(bs.JustifyContent(bs.AlignBetween), bs.AlignItems(bs.AlignCenter), bs.Gap(2), bs.FlexWrap)
		}, []string{`d-flex`, `justify-content-between`, `align-items-center`, `gap-2`, `flex-wrap`, `flex-grow-1`}, nil},
		{"flex column and center", func() hx.INode { return bs.FlexColumn(bs.Center(node())) }, []string{`d-flex flex-column`, `justify-content-center`, `align-items-center`}, nil},
		{"stack and inline", func() hx.INode { return bs.Stack(3, bs.Inline(2, node())) }, []string{`gap-3 vstack`, `gap-2 hstack`}, nil},
		{"grid", func() hx.INode { return bs.Grid(3, 2, node()) }, []string{`d-grid`, `gap-2`, `grid-template-columns:repeat(3, 1fr)`}, nil},
		{"gutters and container", func() hx.INode {
			return bs.ContainerOn(bs.MD, bs.Row().Class(bs.Gutter(3), bs.GutterX(2), bs.GutterY(1)))
		}, []string{`container-md`, `g-3`, `gx-2`, `gy-1`}, nil},
		{"blockquote", func() hx.INode { return bs.Blockquote("Q", "Someone") }, []string{`<figure>`, `<blockquote class="blockquote"><p>Q</p>`, `blockquote-footer">Someone`}, nil},
		{"blockquote without source", func() hx.INode { return bs.Blockquote("Q", "") }, []string{`blockquote`}, []string{`figcaption`}},
		{"lists", func() hx.INode { return hx.Div().Body(bs.ListUnstyled(node()), bs.ListInline(node(), node())) }, []string{`list-unstyled`, `<li><p>x`, `list-inline`, `list-inline-item`}, nil},
		{"stretched link", func() hx.INode { return bs.StretchedLink("/x", "Go") }, []string{`stretched-link`, `href="/x"`}, nil},
		{"table row and caption", func() hx.INode {
			return bs.Table(nil, bs.TableCaption("Cap"), bs.TableRow([]bs.Mod{bs.Danger}, hx.Td()))
		}, []string{`<caption>Cap</caption>`, `table-danger`}, nil},

		{"button check", func() hx.INode { return bs.ButtonCheck("c1", "Toggle", flag) }, []string{`btn-check`, `type="checkbox"`, `checked`, `btn btn-outline-primary`, `for="c1"`}, nil},
		{"button check variant", func() hx.INode { return bs.ButtonCheck("c1", "Toggle", off, bs.Outline(bs.Danger)) }, []string{`btn-outline-danger`}, []string{`checked`, `btn-outline-primary`}},
		{"button radios", func() hx.INode {
			return bs.ButtonRadios("p", text, []bs.Option{{Value: "a", Label: "A"}, {Value: "b", Label: "B"}})
		}, []string{`btn-group`, `role="group"`, `id="p-a"`, `type="radio"`, `for="p-b"`}, nil},
		{"plain text", func() hx.INode { return bs.PlainText("pt", text) }, []string{`form-control-plaintext`, `readonly`, `value="a"`}, nil},
		{"valid", func() hx.INode {
			return hx.Div().Body(bs.IsValid(bs.Input("text", "v", text), flag), bs.ValidFeedback("Good"))
		}, []string{`is-valid`, `valid-feedback">Good`}, nil},
		{"fieldset", func() hx.INode { return bs.Fieldset("Legend", node()) }, []string{`<fieldset>`, `<legend class="col-form-label pt-0">Legend`}, nil},
		{"fieldset without legend", func() hx.INode { return bs.Fieldset("", node()) }, []string{`<fieldset>`}, []string{`legend`}},
		{"inline checks", func() hx.INode {
			return hx.Div().Body(bs.CheckInline("ci", "I", flag), bs.RadiosInline("r", text, bs.Plain("a", "b")))
		}, []string{`form-check form-check-inline`, `id="r-a"`}, nil},
		{"checkbox group", func() hx.INode {
			return bs.CheckboxGroup("Pick", bs.CheckItem{ID: "g1", Label: "One", Value: flag}, bs.CheckItem{ID: "g2", Label: "Two", Value: off})
		}, []string{`<fieldset>`, `Pick`, `id="g1"`, `id="g2"`}, nil},
		{"color and datalist", func() hx.INode {
			return hx.Div().Body(bs.ColorInput("col", hx.Signal("#ff0000")), bs.Datalist("dl", "a", "b"))
		}, []string{`form-control form-control-color`, `type="color"`, `<datalist id="dl">`, `<option value="a">`}, nil},
		{"input group size and check", func() hx.INode {
			return bs.InputGroupSize(bs.Sm, bs.InputGroupCheck("on", flag), bs.Input("text", "ig", text))
		}, []string{`input-group input-group-sm`, `input-group-text`, `form-check-input mt-0`, `aria-label="on"`}, nil},
		{"form row", func() hx.INode { return bs.FormRow(bs.ColOn(bs.MD, 6)) }, []string{`class="g-3 row"`, `col-md-6`}, nil},

		{"alert parts", func() hx.INode { return bs.Alert(hx.Div().Body(bs.AlertHeading("H"), bs.AlertLink("/a", "link"))) }, []string{`alert-heading">H`, `alert-link`}, nil},
		{"card group", func() hx.INode { return bs.CardGroup(bs.Card()) }, []string{`card-group`}, nil},
		{"card grid", func() hx.INode { return bs.CardGrid(3, bs.Card(), bs.Card()) }, []string{`row-cols-1`, `row-cols-md-3`, `g-4`, `class="col"`}, nil},
		{"card overlay", func() hx.INode { return bs.CardImageOverlay(node()) }, []string{`card-img-overlay`}, nil},
		{"card header tabs", func() hx.INode { return bs.Card(bs.CardHeaderTabs(tab, []bs.Tab{{ID: "x", Label: "X"}})) }, []string{`card-header`, `card-header-tabs`, `nav-tabs`}, nil},
		{"card header pills", func() hx.INode { return bs.CardHeaderTabs(tab, []bs.Tab{{ID: "x", Label: "X"}}, bs.Pills) }, []string{`card-header-pills`, `nav-pills`}, []string{`card-header-tabs`}},
		{"badge on a button", func() hx.INode { return bs.WithBadge(bs.Button("Inbox"), "9") }, []string{`position-relative`, `d-inline-block`, `translate-middle`, `text-bg-danger`, `>9</span>`}, nil},
		{"badge variant", func() hx.INode { return bs.WithBadge(bs.Button("Inbox"), "9", bs.Success) }, []string{`text-bg-success`}, []string{`text-bg-danger`}},
		{"loading button idle", func() hx.INode { return bs.LoadingButton("Save", off) }, []string{`btn btn-primary`, `>Save</span>`}, []string{`spinner-border`, `disabled`}},
		{"loading button busy", func() hx.INode { return bs.LoadingButton("Save", flag, bs.Success) }, []string{`btn-success`, `spinner-border-sm`, `disabled`}, nil},
		{"stacked progress", func() hx.INode {
			return bs.ProgressStacked(bs.ProgressSegment(15, 100), bs.ProgressSegment(30, 100, bs.Success), bs.ProgressSegmentOf(n, 100, bs.Info))
		}, []string{`progress-stacked`, `width:15%`, `width:30%`, `width:20%`, `text-bg-success`, `aria-valuenow="20"`}, nil},

		{"dropup end", func() hx.INode { return bs.Dropdown(off, "Menu", items, bs.Dropup, bs.MenuEnd, bs.Primary) }, []string{`class="dropup"`, `dropdown-menu-end`, `btn-primary`}, []string{`class="dropdown"`}},
		{"dropdown default", func() hx.INode { return bs.Dropdown(off, "Menu", items) }, []string{`class="dropdown"`, `btn-secondary`, `data-bs-popper="static"`}, []string{`dropdown-menu-end`}},
		{"nav dropdown", func() hx.INode { return bs.Nav(bs.NavDropdown(flag, "More", items)) }, []string{`<li class="dropdown nav-item">`, `dropdown-toggle nav-link`, `dropdown-menu show`, `aria-expanded="true"`}, nil},
		{"split dropdown", func() hx.INode {
			return bs.SplitDropdown(bs.SplitDropdownOptions{Open: off, Label: "Action", OnClick: func() {}, Items: items, Mods: []bs.Mod{bs.Danger, bs.Dropend}})
		}, []string{`btn-group dropend`, `btn btn-danger`, `dropdown-toggle-split`, `Toggle Dropdown`}, nil},
		{"modal scrollable", func() hx.INode {
			return bs.Modal(bs.ModalOptions{Open: flag, Title: "T", Body: node, Scrollable: true, Size: bs.Lg})
		}, []string{`modal-dialog-scrollable`, `modal-lg`}, nil},
		{"modal fullscreen", func() hx.INode {
			return bs.Modal(bs.ModalOptions{Open: flag, Title: "T", Body: node, Fullscreen: true, Centered: true})
		}, []string{`modal-fullscreen`, `modal-dialog-centered`}, nil},
		{"navbar text", func() hx.INode { return bs.Navbar("lg", bs.NavbarText("Hi")).Class(bs.FixedTop) }, []string{`navbar-text">Hi`, `fixed-top`}, nil},
	}
	for _, c := range cases {
		got := render(c.build)
		for _, w := range c.want {
			if !strings.Contains(got, w) {
				t.Errorf("%s: missing %q in\n%s", c.name, w, got)
			}
		}
		for _, w := range c.not {
			if strings.Contains(got, w) {
				t.Errorf("%s: has %q in\n%s", c.name, w, got)
			}
		}
	}
}
