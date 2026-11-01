package bs

import (
	"strconv"
	"strings"

	"github.com/deltegui/hx"
)

// A closed overlay is not in the tree: the content of a Modal, an Offcanvas and
// a Toast is built each time it opens, so a closed one costs nothing. Bootstrap
// shows them by adding the class "show" with JavaScript; here the node is
// created already shown.

func build(f func() hx.INode) hx.INode {
	if f == nil {
		return nil
	}
	return f()
}

// shown makes the node of an overlay only while open is true. mk is called
// without tracking, so what it reads does not rebuild the overlay.
func shown(open hx.Gettable[bool], mk func() hx.INode) hx.INode {
	return hx.Dynamic(func() hx.INode {
		if !open.Get() {
			return nil
		}
		var n hx.INode
		hx.Untrack(func() { n = mk() })
		return n
	})
}

func backdrop(class string, onClick func()) hx.INode {
	b := hx.Div().Class(class, "fade", "show")
	if onClick != nil {
		b.OnClick(func(hx.EventContext) { onClick() })
	}
	return b
}

// Modal is a dialog over the page, with a background that closes it. body and
// footer are built each time it opens (footer may be nil). onClose runs when the
// user clicks the background or the close button. Mods: Sm, Lg, Xl, and
// Mod("fullscreen"). For one in the middle of the screen use ModalCentered.
func Modal(open hx.Gettable[bool], title string, onClose func(), body, footer func() hx.INode, mods ...Mod) hx.INode {
	return modal(open, title, onClose, body, footer, false, mods)
}

// ModalCentered is a Modal in the middle of the screen, not near the top.
func ModalCentered(open hx.Gettable[bool], title string, onClose func(), body, footer func() hx.INode, mods ...Mod) hx.INode {
	return modal(open, title, onClose, body, footer, true, mods)
}

func modal(open hx.Gettable[bool], title string, onClose func(), body, footer func() hx.INode, centered bool, mods []Mod) hx.INode {
	return shown(open, func() hx.INode {
		parts := []hx.INode{
			hx.Div().Class("modal-header").Body(
				hx.H5().Class("modal-title").Text(title),
				CloseButton(onClose)),
			hx.Div().Class("modal-body").Body(build(body)),
		}
		if footer != nil {
			parts = append(parts, hx.Div().Class("modal-footer").Body(build(footer)))
		}
		content := hx.Div().Class("modal-content").BodyList(parts)
		dialog := with(hx.Div().Class("modal-dialog").Body(content), "modal", mods)
		dialog.ClassIf("modal-dialog-centered", centered)
		// The dialog and the content are above the background, so a click on
		// the shadow reaches it: the .modal does not take the pointer.
		modal := hx.Div().Class("modal", "fade", "show", "d-block").
			Attribute("tabindex", "-1").Attribute("role", "dialog").Attribute("aria-modal", "true").
			Style("pointer-events", "none").Body(dialog)
		return hx.Noop().Body(backdrop("modal-backdrop", onClose), modal)
	})
}

// Offcanvas is a panel that slides from a side. side is "start", "end", "top" or
// "bottom". body is built each time it opens.
func Offcanvas(open hx.Gettable[bool], side, title string, onClose func(), body func() hx.INode) hx.INode {
	return shown(open, func() hx.INode {
		panel := hx.Div().Class("offcanvas", "offcanvas-"+side, "show").
			Attribute("tabindex", "-1").Attribute("role", "dialog").Body(
			hx.Div().Class("offcanvas-header").Body(
				hx.H5().Class("offcanvas-title").Text(title),
				CloseButton(onClose)),
			hx.Div().Class("offcanvas-body").Body(build(body)))
		return hx.Noop().Body(backdrop("offcanvas-backdrop", onClose), panel)
	})
}

// Dropdown is a button that opens a menu. open says if it is open and the button
// toggles it; items go in the menu (DropdownItem, DropdownDivider). A click
// outside closes it. Choosing an item does not close it: do open.Set(false) in
// the handler of the item. Mods: the colors and sizes of the button.
func Dropdown(open *hx.SignalT[bool], label string, items []hx.INode, mods ...Mod) hx.INode {
	toggle := hx.Button().Type("button").Class("btn", "dropdown-toggle").Attribute("aria-haspopup", "true").Text(label)
	if len(mods) == 0 {
		toggle.Class("btn-secondary")
	} else {
		with(toggle, "btn", mods)
	}
	toggle.OnClick(func(hx.EventContext) { open.Set(!open.Get()) })
	toggle.Effect(func() { toggle.Attribute("aria-expanded", strconv.FormatBool(open.Get())) })
	menu := hx.Ul().Class("dropdown-menu").Body(items...)
	BindClass(menu, "show", open)
	// An invisible layer under the menu catches the click that closes it.
	away := hx.Dynamic(func() hx.INode {
		if !open.Get() {
			return nil
		}
		return hx.Div().Attribute("aria-hidden", "true").
			Style("position", "fixed").Style("inset", "0").Style("z-index", "999").
			OnClick(func(hx.EventContext) { open.Set(false) })
	})
	menu.Style("z-index", "1000")
	return hx.Div().Class("dropdown").Body(toggle, away, menu)
}

// DropdownItem is an item of a Dropdown; with href empty it is a button. Chain
// "active" or "disabled" on it.
func DropdownItem(label, href string, onClick func()) hx.INode {
	var a hx.INode
	if href != "" {
		a = hx.A().Href(href)
	} else {
		a = hx.Button().Type("button")
	}
	a.Class("dropdown-item").Text(label)
	if onClick != nil {
		a.OnClick(func(hx.EventContext) { onClick() })
	}
	return hx.Li().Body(a)
}

func DropdownDivider() hx.INode { return hx.Li().Body(hx.Hr().Class("dropdown-divider")) }

func DropdownHeader(text string) hx.INode {
	return hx.Li().Body(hx.H6().Class("dropdown-header").Text(text))
}

// Collapse is content that folds: it shows while open is true. Pair it with a
// button that toggles the signal. The content is built once and hidden with CSS,
// so what it holds (a form) keeps its state.
func Collapse(open hx.Gettable[bool], children ...hx.INode) hx.INode {
	return BindClass(hx.Div().Class("collapse").Body(children...), "show", open)
}

// AccordionItem is a section of an Accordion.
type AccordionItem struct {
	ID, Title string
	Body      func() hx.INode
}

// Accordion is a list of sections of which one is open. active holds the ID of
// the open one ("" for none) and clicking a title opens it or closes it. The
// body of a section is built once, so what it holds keeps its state.
func Accordion(active *hx.SignalT[string], items ...AccordionItem) hx.INode {
	nodes := make([]hx.INode, len(items))
	for i, it := range items {
		button := hx.Button().Type("button").Class("accordion-button").Text(it.Title)
		button.OnClick(func(hx.EventContext) {
			if active.Get() == it.ID {
				active.Set("")
			} else {
				active.Set(it.ID)
			}
		})
		button.Effect(func() {
			on := active.Get() == it.ID
			button.ClassIf("collapsed", !on)
			button.Attribute("aria-expanded", strconv.FormatBool(on))
		})
		panel := hx.Div().Class("accordion-collapse", "collapse").Body(hx.Div().Class("accordion-body").Body(build(it.Body)))
		panel.Effect(func() { panel.ClassIf("show", active.Get() == it.ID) })
		nodes[i] = hx.Div().Class("accordion-item").Body(hx.H2().Class("accordion-header").Body(button), panel)
	}
	return hx.Div().Class("accordion").Body(nodes...)
}

// Toast is a short message in a corner. It is in the tree while open is true;
// body is built each time. Put it in a ToastContainer.
func Toast(open hx.Gettable[bool], title string, onClose func(), body func() hx.INode) hx.INode {
	return shown(open, func() hx.INode {
		return hx.Div().Class("toast", "show").Attribute("role", "alert").Attribute("aria-live", "assertive").Body(
			hx.Div().Class("toast-header").Body(
				hx.Strong().Class("me-auto").Text(title),
				CloseButton(onClose)),
			hx.Div().Class("toast-body").Body(build(body)))
	})
}

// ToastContainer holds the toasts in a corner of the screen. position is a
// Bootstrap position: "top-0 end-0", "bottom-0 start-0"...
func ToastContainer(position string, children ...hx.INode) hx.INode {
	c := hx.Div().Class("toast-container", "position-fixed", "p-3").Body(children...)
	for _, p := range strings.Fields(position) {
		c.Class(p)
	}
	return c
}
