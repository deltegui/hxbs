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

// closer is the function that closes something: onClose, or, when it is nil and
// open is a signal, a function that sets it to false.
func closer(open hx.Gettable[bool], onClose func()) func() {
	if onClose != nil {
		return onClose
	}
	if s, ok := open.(hx.Settable[bool]); ok {
		return func() { s.Set(false) }
	}
	return nil
}

// ModalOptions is what a Modal needs. Only Open is required.
type ModalOptions struct {
	// Open says if the modal is in the tree. If it is a signal, the modal closes
	// itself with it (see OnClose).
	Open  hx.Gettable[bool]
	Title string
	// Body and Footer are built each time the modal opens. Footer may be nil.
	Body, Footer func() hx.INode
	// OnClose runs when the user clicks the background or the close button. By
	// default it sets Open to false, when Open is a signal.
	OnClose func()
	// Size is Sm, Lg or Xl (the default is the medium one).
	Size Mod
	// Centered puts it in the middle of the screen, not near the top. Scrollable
	// makes the body scroll instead of the page. Fullscreen fills the window.
	Centered, Scrollable, Fullscreen bool
}

// Modal is a dialog over the page, with a background that closes it:
//
//	bs.Modal(bs.ModalOptions{Open: open, Title: "Edit", Body: form, Footer: buttons, Size: bs.Lg})
func Modal(o ModalOptions) hx.INode {
	onClose := closer(o.Open, o.OnClose)
	return shown(o.Open, func() hx.INode {
		parts := []hx.INode{
			hx.Div().Class("modal-header").Body(
				hx.H5().Class("modal-title").Text(o.Title),
				CloseButton(onClose)),
			hx.Div().Class("modal-body").Body(build(o.Body)),
		}
		if o.Footer != nil {
			parts = append(parts, hx.Div().Class("modal-footer").Body(build(o.Footer)))
		}
		content := hx.Div().Class("modal-content").BodyList(parts)
		dialog := hx.Div().Class("modal-dialog").Body(content)
		if o.Size != "" {
			dialog.Class("modal-" + string(o.Size))
		}
		dialog.ClassIf("modal-dialog-centered", o.Centered)
		dialog.ClassIf("modal-dialog-scrollable", o.Scrollable)
		dialog.ClassIf("modal-fullscreen", o.Fullscreen)
		// The dialog and the content are above the background, so a click on
		// the shadow reaches it: the .modal does not take the pointer.
		modal := hx.Div().Class("modal", "fade", "show", "d-block").
			Attribute("tabindex", "-1").Attribute("role", "dialog").Attribute("aria-modal", "true").
			Style("pointer-events", "none").Body(dialog)
		return hx.Noop().Body(backdrop("modal-backdrop", onClose), modal)
	})
}

// OffcanvasOptions is what an Offcanvas needs. Only Open is required.
type OffcanvasOptions struct {
	// Open says if the panel is in the tree. If it is a signal, the panel closes
	// itself with it (see OnClose).
	Open  hx.Gettable[bool]
	Title string
	// Side is "start" (the default), "end", "top" or "bottom".
	Side string
	// Body is built each time the panel opens.
	Body func() hx.INode
	// OnClose runs when the user clicks the background or the close button. By
	// default it sets Open to false, when Open is a signal.
	OnClose func()
}

// Offcanvas is a panel that slides from a side:
//
//	bs.Offcanvas(bs.OffcanvasOptions{Open: open, Title: "Menu", Side: "end", Body: menu})
func Offcanvas(o OffcanvasOptions) hx.INode {
	onClose := closer(o.Open, o.OnClose)
	side := o.Side
	if side == "" {
		side = "start"
	}
	return shown(o.Open, func() hx.INode {
		panel := hx.Div().Class("offcanvas", "offcanvas-"+side, "show").
			Attribute("tabindex", "-1").Attribute("role", "dialog").Body(
			hx.Div().Class("offcanvas-header").Body(
				hx.H5().Class("offcanvas-title").Text(o.Title),
				CloseButton(onClose)),
			hx.Div().Class("offcanvas-body").Body(build(o.Body)))
		return hx.Noop().Body(backdrop("offcanvas-backdrop", onClose), panel)
	})
}

// Dropdown is a button that opens a menu. open says if it is open and the button
// toggles it; items go in the menu (DropdownItem, DropdownDivider). A click
// outside closes it. Choosing an item does not close it: do open.Set(false) in
// the handler of the item. Mods: the colors and sizes of the button.
// Where a dropdown opens and where its menu is aligned. They go in the mods of
// Dropdown, NavDropdown and SplitDropdown; any other mod is the variant of the button.
const (
	Dropup    Mod = "dropup"
	Dropend   Mod = "dropend"
	Dropstart Mod = "dropstart"
	MenuEnd   Mod = "menu-end"
)

func dropMods(mods []Mod) (button []Mod, root string, menuEnd bool) {
	root = "dropdown"
	for _, m := range mods {
		switch m {
		case Dropup, Dropend, Dropstart:
			root = string(m)
		case MenuEnd:
			menuEnd = true
		default:
			button = append(button, m)
		}
	}
	return
}

func styleButton(b hx.INode, mods []Mod) {
	if len(mods) == 0 {
		b.Class("btn-secondary")
	} else {
		with(b, "btn", mods)
	}
}

func wireToggle(toggle hx.INode, open *hx.SignalT[bool]) {
	toggle.OnClick(func(hx.EventContext) { open.Set(!open.Get()) })
	toggle.Effect(func() { toggle.Attribute("aria-expanded", strconv.FormatBool(open.Get())) })
}

// dropdownMenu is the menu and the invisible layer that closes it.
func dropdownMenu(open *hx.SignalT[bool], items []hx.INode, menuEnd bool) (away, menu hx.INode) {
	// Bootstrap places the menu (top, left/right, dropup, dropend, menu-end) only when
	// it has data-bs-popper; its JavaScript sets it, and "static" is the value it uses
	// when there is no Popper.
	menu = hx.Ul().Class("dropdown-menu").Attribute("data-bs-popper", "static").Body(items...)
	menu.ClassIf("dropdown-menu-end", menuEnd)
	BindClass(menu, "show", open)
	away = hx.Dynamic(func() hx.INode {
		if !open.Get() {
			return nil
		}
		return hx.Div().Attribute("aria-hidden", "true").
			Style("position", "fixed").Style("inset", "0").Style("z-index", "999").
			OnClick(func(hx.EventContext) { open.Set(false) })
	})
	menu.Style("z-index", "1000")
	return
}

func Dropdown(open *hx.SignalT[bool], label string, items []hx.INode, mods ...Mod) hx.INode {
	button, root, end := dropMods(mods)
	toggle := hx.Button().Type("button").Class("btn", "dropdown-toggle").Attribute("aria-haspopup", "true").Text(label)
	styleButton(toggle, button)
	wireToggle(toggle, open)
	away, menu := dropdownMenu(open, items, end)
	return hx.Div().Class(root).Body(toggle, away, menu)
}

// NavDropdown is a dropdown that is an item of a Nav or a NavbarNav.
func NavDropdown(open *hx.SignalT[bool], label string, items []hx.INode, mods ...Mod) hx.INode {
	_, root, end := dropMods(mods)
	toggle := hx.Button().Type("button").Class("nav-link", "dropdown-toggle").Attribute("aria-haspopup", "true").Text(label)
	wireToggle(toggle, open)
	away, menu := dropdownMenu(open, items, end)
	return hx.Li().Class("nav-item", root).Body(toggle, away, menu)
}

// SplitDropdownOptions is what a SplitDropdown needs.
type SplitDropdownOptions struct {
	// Open says if the menu is open; the arrow toggles it.
	Open  *hx.SignalT[bool]
	Label string
	// OnClick is the action of the main button. It may be nil.
	OnClick func()
	// Items go in the menu (DropdownItem, DropdownDivider).
	Items []hx.INode
	// Mods are the variant of the button (bs.Danger) and where the menu opens
	// (bs.Dropup, bs.Dropend, bs.Dropstart, bs.MenuEnd).
	Mods []Mod
}

// SplitDropdown is a button with its own action and, next to it, the arrow that opens the menu.
func SplitDropdown(o SplitDropdownOptions) hx.INode {
	button, root, end := dropMods(o.Mods)
	main := hx.Button().Type("button").Class("btn").Text(o.Label)
	styleButton(main, button)
	if o.OnClick != nil {
		main.OnClick(func(hx.EventContext) { o.OnClick() })
	}
	toggle := hx.Button().Type("button").Class("btn", "dropdown-toggle", "dropdown-toggle-split").
		Attribute("aria-haspopup", "true").Body(hx.Span().Class("visually-hidden").Text("Toggle Dropdown"))
	styleButton(toggle, button)
	wireToggle(toggle, o.Open)
	away, menu := dropdownMenu(o.Open, o.Items, end)
	group := hx.Div().Class("btn-group").Body(main, toggle, away, menu)
	return group.ClassIf(root, root != "dropdown")
}

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

// ToastOptions is what a Toast needs. Only Open is required.
type ToastOptions struct {
	// Open says if the toast is in the tree. If it is a signal, the toast closes
	// itself with it (see OnClose).
	Open  hx.Gettable[bool]
	Title string
	// Body is built each time the toast shows.
	Body func() hx.INode
	// OnClose runs when the user clicks the close button. By default it sets Open
	// to false, when Open is a signal.
	OnClose func()
}

// Toast is a short message in a corner. Put it in a ToastContainer.
func Toast(o ToastOptions) hx.INode {
	onClose := closer(o.Open, o.OnClose)
	return shown(o.Open, func() hx.INode {
		return hx.Div().Class("toast", "show").Attribute("role", "alert").Attribute("aria-live", "assertive").Body(
			hx.Div().Class("toast-header").Body(
				hx.Strong().Class("me-auto").Text(o.Title),
				CloseButton(onClose)),
			hx.Div().Class("toast-body").Body(build(o.Body)))
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
