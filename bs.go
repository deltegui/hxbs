// Package bs wraps Bootstrap 5.3 in Hix nodes: every component of the framework
// is a function that returns an hx.INode with the right tags and classes.
//
// It uses only Bootstrap's CSS. What Bootstrap does with JavaScript (modal,
// offcanvas, dropdown, collapse, accordion, tabs, toast, carousel) is here a
// Signal, so the same code runs on the server (hx/ssr), in the browser (hx/wasm)
// and when the browser hydrates what the server wrote. It imports only hx.
//
// A component takes what it needs to exist (its text, its children, the signal
// it reads) and returns a node. The variants of a component are Mods, which the
// component prefixes with its own name, so Primary is "btn-primary" on a Button
// and "alert-primary" on an Alert:
//
//	bs.Button("Save", bs.Primary, bs.Lg)
//	bs.Alert(hx.Span().Text("Saved"), bs.Success)
//	bs.Button("Delete", bs.Outline(bs.Danger))
//
// Anything else (utilities such as "mb-3" or "d-flex") is a class, chained as
// usual: bs.Card(...).Class("mb-3").
//
// Bootstrap is a stylesheet that the page loads:
//
//	<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css">
//
// (the constant CSS has that address).
package bs

import (
	"strconv"
	"strings"

	"github.com/deltegui/hx"
)

// CSS is the stylesheet of the version of Bootstrap that this package is written for.
const CSS = "https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css"

// Icons is the stylesheet of Bootstrap Icons, which is not part of Bootstrap:
// Icon("house") takes the name of an icon of that font.
const Icons = "https://cdn.jsdelivr.net/npm/bootstrap-icons@1.13.1/font/bootstrap-icons.min.css"

// A Mod is a variant of a component: a color, a size or a style. The component
// that receives it puts its own prefix in front, so Mod("primary") is "btn-primary"
// in a Button and "bg-primary" nowhere (see BgColor for utilities).
type Mod string

// Colors.
const (
	Primary   Mod = "primary"
	Secondary Mod = "secondary"
	Success   Mod = "success"
	Danger    Mod = "danger"
	Warning   Mod = "warning"
	Info      Mod = "info"
	Light     Mod = "light"
	Dark      Mod = "dark"
	Link      Mod = "link"
)

// Sizes (buttons, inputs, pagination, spinners, list groups...).
const (
	Sm Mod = "sm"
	Lg Mod = "lg"
	Xl Mod = "xl"
)

// Shapes and states of a table, a list group, a button group or a spinner.
const (
	Striped    Mod = "striped"
	Hover      Mod = "hover"
	Bordered   Mod = "bordered"
	Borderless Mod = "borderless"
	Flush      Mod = "flush"
	Numbered   Mod = "numbered"
	Horizontal Mod = "horizontal"
	Vertical   Mod = "vertical"
	Grow       Mod = "grow"
	Animated   Mod = "animated"
)

// Outline is the outlined variant of a color: Outline(Danger) on a Button is
// "btn-outline-danger".
func Outline(color Mod) Mod { return "outline-" + color }

// classes turns mods into classes with the prefix of the component.
func classes(prefix string, mods []Mod) []string {
	out := make([]string, len(mods))
	for i, m := range mods {
		out[i] = prefix + "-" + string(m)
	}
	return out
}

// with adds the mods of a component to its node.
func with(n hx.INode, prefix string, mods []Mod) hx.INode {
	if len(mods) > 0 {
		n.Class(classes(prefix, mods)...)
	}
	return n
}

// BindClass adds the class to the node while the signal is true, and removes it
// otherwise. It is how a component reacts to a Signal without being rebuilt:
//
//	bs.BindClass(bs.Button("Save", bs.Primary), "disabled", saving)
func BindClass(n hx.INode, class string, on hx.Gettable[bool]) hx.INode {
	return n.Effect(func() { n.ClassIf(class, on.Get()) })
}

// Breakpoints, for ColOn and RowColsOn.
const (
	SM  = "sm"
	MD  = "md"
	LG  = "lg"
	XL  = "xl"
	XXL = "xxl"
)

// Color utilities: TextColor(Primary) is "text-primary", BgColor(Light) is
// "bg-light", BorderColor(Danger) is "border-danger".
func TextColor(c Mod) string   { return "text-" + string(c) }
func BgColor(c Mod) string     { return "bg-" + string(c) }
func BorderColor(c Mod) string { return "border-" + string(c) }

// TextBg is the "text-bg-primary" of Badge and Card: a background with the text
// color that reads well on it.
func TextBg(c Mod) string { return "text-bg-" + string(c) }

// Space is a spacing utility. property is "m" or "p", side is "t", "b", "s",
// "e", "x", "y" or "" for all, and n is 0 to 5: Space("m", "t", 3) is "mt-3".
func Space(property, side string, n int) string {
	return property + side + "-" + strconv.Itoa(n)
}

// Gap is the gap of a flex or grid container: Gap(3) is "gap-3".
func Gap(n int) string { return "gap-" + strconv.Itoa(n) }

// Icon is an icon of Bootstrap Icons (load Icons in the page): Icon("house").
func Icon(name string) hx.INode {
	return hx.I().Class("bi", "bi-"+strings.TrimPrefix(name, "bi-")).Attribute("aria-hidden", "true")
}

func button(label string, kind string, mods []Mod) hx.INode {
	b := hx.Button().Type(kind).Class("btn").Text(label)
	if len(mods) == 0 {
		return b.Class("btn-primary")
	}
	return with(b, "btn", mods)
}
