package bs

import (
	"strings"

	"github.com/deltegui/hx"
)

// Option is an entry of a Select or a group of Radios.
type Option struct{ Value, Label string }

// Plain makes the Options of values whose label is the value: Plain("a", "b").
func Plain(values ...string) []Option {
	out := make([]Option, len(values))
	for i, v := range values {
		out[i] = Option{v, v}
	}
	return out
}

// Label is the label of a control; id is the id of that control. Bootstrap
// has the label above the control: for a checkbox use Checkbox.
func Label(id, text string) hx.INode {
	return hx.Label().Class("form-label").Attribute("for", id).Text(text)
}

// FormText is the small gray text under a control.
func FormText(text string) hx.INode { return hx.Div().Class("form-text").Text(text) }

// Input is an input of the given type ("text", "email", "password", "number",
// "date", "tel", "url", "search"...), bound to the signal. Mods: Sm, Lg.
func Input(kind, id string, value *hx.SignalT[string], mods ...Mod) hx.INode {
	in := hx.Input().Type(kind).Id(id).Class("form-control").BindValue(value).BindOnInput(value)
	return with(in, "form-control", mods)
}

// TextArea is a multi line text.
func TextArea(id string, value *hx.SignalT[string], mods ...Mod) hx.INode {
	t := hx.TextArea().Id(id).Class("form-control").BindValue(value).BindOnInput(value)
	return with(t, "form-control", mods)
}

// Select is a dropdown list. It has an empty first option with the text of
// blank (none if blank is empty). Mods: Sm, Lg.
func Select(id string, value *hx.SignalT[string], blank string, options []Option, mods ...Mod) hx.INode {
	sel := hx.Select().Id(id).Class("form-select")
	var nodes []hx.INode
	if blank != "" {
		nodes = append(nodes, hx.Option().Value("").Text(blank))
	}
	for _, o := range options {
		nodes = append(nodes, hx.Option().Value(o.Value).Text(o.Label))
	}
	sel.BodyList(nodes).BindValue(value).BindOnChange(value)
	return with(sel, "form-select", mods)
}

// FileInput is a file picker; onChange runs when the user chooses (read the files
// with hxhttp.FilesOf(ctx.Target) in the browser).
func FileInput(id string, onChange func(hx.EventContext)) hx.INode {
	in := hx.Input().Type("file").Id(id).Attribute("name", "file").Class("form-control")
	if onChange != nil {
		in.OnChange(onChange)
	}
	return in
}

// Range is a slider; its value is the text of the number.
func Range(id string, value *hx.SignalT[string]) hx.INode {
	return hx.Input().Type("range").Id(id).Class("form-range").BindValue(value).BindOnInput(value)
}

// Checkbox is a checkbox with its label, bound to the signal. id links the
// label to the input, so it must be unique in the page.
func Checkbox(id, label string, value *hx.SignalT[bool]) hx.INode {
	return hx.Div().Class("form-check").Body(
		hx.Input().Type("checkbox").Id(id).Class("form-check-input").BindChecked(value).BindOnChecked(value),
		hx.Label().Class("form-check-label").Attribute("for", id).Text(label))
}

// Switch is a Checkbox that looks like a toggle.
func Switch(id, label string, value *hx.SignalT[bool]) hx.INode {
	return Checkbox(id, label, value).Class("form-switch")
}

// Radio is one radio button: it is checked while the signal has its value, and
// it sets the signal when chosen. Radios of a group share name and signal.
func Radio(name, value, label string, selected *hx.SignalT[string]) hx.INode {
	id := name + "-" + value
	input := hx.Input().Type("radio").Id(id).Class("form-check-input").Attribute("name", name).Value(value)
	input.Effect(func() { input.Checked(selected.Get() == value) })
	input.OnChange(func(hx.EventContext) { selected.Set(value) })
	return hx.Div().Class("form-check").Body(input,
		hx.Label().Class("form-check-label").Attribute("for", id).Text(label))
}

// Radios is a Radio for each option. Chain "form-check-inline" on each, or use
// Radio yourself, to put them in a line.
func Radios(name string, selected *hx.SignalT[string], options []Option) hx.INode {
	nodes := make([]hx.INode, len(options))
	for i, o := range options {
		nodes[i] = Radio(name, o.Value, o.Label, selected)
	}
	return hx.Div().Body(nodes...)
}

// FormField is the common case: a label, a control and the errors of that field
// as a red text (nothing if errs is nil or empty). The control gets "is-invalid"
// while there are errors and the text follows the signal. id is the id of the
// control, as it was given to Input, Select or TextArea.
func FormField(label, id string, control hx.INode, errs hx.Gettable[[]string]) hx.INode {
	f := hx.Div().Class("mb-3")
	if label != "" {
		f.Body(Label(id, label), control, feedback(control, errs))
		return f
	}
	return f.Body(control, feedback(control, errs))
}

func feedback(control hx.INode, errs hx.Gettable[[]string]) hx.INode {
	fb := hx.Div().Class("invalid-feedback")
	if errs != nil {
		fb.Effect(func() {
			e := strings.Join(errs.Get(), " · ")
			fb.Text(e)
			control.ClassIf("is-invalid", e != "")
		})
	}
	return fb
}

// FloatingLabel is a control with its label inside it. The control needs a
// placeholder to work (Bootstrap uses it to know if it is empty): Placeholder
// is set to the label if the control has none.
func FloatingLabel(label, id string, control hx.INode) hx.INode {
	if control.GetAttribute("placeholder") == "" {
		control.Placeholder(label)
	}
	return hx.Div().Class("form-floating").Body(control, hx.Label().Attribute("for", id).Text(label))
}

// InputGroup sticks controls, buttons and InputGroupText together on one line.
func InputGroup(children ...hx.INode) hx.INode {
	return hx.Div().Class("input-group").Body(children...)
}

// InputGroupText is the gray box with a text or an icon of an InputGroup.
func InputGroupText(children ...hx.INode) hx.INode {
	return hx.Span().Class("input-group-text").Body(children...)
}

// ButtonCheck is a checkbox that looks like a button (btn-check and a label.btn).
// Without a variant it is btn-outline-primary.
func ButtonCheck(id, label string, value *hx.SignalT[bool], mods ...Mod) hx.INode {
	in := hx.Input().Type("checkbox").Id(id).Class("btn-check").Attribute("autocomplete", "off").
		BindChecked(value).BindOnChecked(value)
	return hx.Noop().Body(in, buttonLabel(id, label, mods))
}

// ButtonRadios is a group of radios that look like buttons.
func ButtonRadios(name string, selected *hx.SignalT[string], options []Option, mods ...Mod) hx.INode {
	nodes := make([]hx.INode, 0, 2*len(options))
	for _, o := range options {
		id := name + "-" + o.Value
		in := hx.Input().Type("radio").Id(id).Class("btn-check").Attribute("name", name).
			Attribute("autocomplete", "off").Value(o.Value)
		in.Effect(func() { in.Checked(selected.Get() == o.Value) })
		in.OnChange(func(hx.EventContext) { selected.Set(o.Value) })
		nodes = append(nodes, in, buttonLabel(id, o.Label, mods))
	}
	return hx.Div().Class("btn-group").Attribute("role", "group").BodyList(nodes)
}

func buttonLabel(id, label string, mods []Mod) hx.INode {
	l := hx.Label().Class("btn").Attribute("for", id).Text(label)
	if len(mods) == 0 {
		return l.Class("btn-outline-primary")
	}
	return with(l, "btn", mods)
}

// PlainText shows a value as read-only text with the look of a form row.
func PlainText(id string, value hx.Gettable[string]) hx.INode {
	return hx.Input().Type("text").Id(id).Class("form-control-plaintext").
		Attribute("readonly", "").BindValue(value)
}

// ValidFeedback is the message of a control marked with IsValid. Bootstrap shows it
// only when it is after the control (FormField puts the invalid one there too).
func ValidFeedback(text string) hx.INode { return hx.Div().Class("valid-feedback").Text(text) }

// IsValid marks the control with is-valid while on is true.
func IsValid(control hx.INode, on hx.Gettable[bool]) hx.INode {
	return BindClass(control, "is-valid", on)
}

// Fieldset groups controls under a legend.
func Fieldset(legend string, children ...hx.INode) hx.INode {
	if legend != "" {
		children = append([]hx.INode{hx.Legend().Class("col-form-label", "pt-0").Text(legend)}, children...)
	}
	return hx.Fieldset().Body(children...)
}

// CheckInline is a Checkbox that shares its line with the next one.
func CheckInline(id, label string, value *hx.SignalT[bool]) hx.INode {
	return Checkbox(id, label, value).Class("form-check-inline")
}

// RadiosInline is Radios in one line.
func RadiosInline(name string, selected *hx.SignalT[string], options []Option) hx.INode {
	nodes := make([]hx.INode, len(options))
	for i, o := range options {
		nodes[i] = Radio(name, o.Value, o.Label, selected).Class("form-check-inline")
	}
	return hx.Div().Body(nodes...)
}

// CheckItem is one checkbox of a CheckboxGroup.
type CheckItem struct {
	ID, Label string
	Value     *hx.SignalT[bool]
}

// CheckboxGroup is a Fieldset of checkboxes.
func CheckboxGroup(legend string, items ...CheckItem) hx.INode {
	nodes := make([]hx.INode, len(items))
	for i, it := range items {
		nodes[i] = Checkbox(it.ID, it.Label, it.Value)
	}
	return Fieldset(legend, nodes...)
}

// ColorInput is <input type=color>, bound to a "#rrggbb" text.
func ColorInput(id string, value *hx.SignalT[string]) hx.INode {
	return hx.Input().Type("color").Id(id).Class("form-control", "form-control-color").
		BindValue(value).BindOnInput(value)
}

// Datalist is the list of suggestions of an Input: input.Attribute("list", id).
func Datalist(id string, values ...string) hx.INode {
	opts := make([]hx.INode, len(values))
	for i, v := range values {
		opts[i] = hx.Option().Value(v)
	}
	return hx.Datalist().Id(id).BodyList(opts)
}

// InputGroupSize is an InputGroup of a size (bs.Sm, bs.Lg).
func InputGroupSize(size Mod, children ...hx.INode) hx.INode {
	return InputGroup(children...).Class("input-group-" + string(size))
}

// InputGroupCheck is a checkbox inside an input group.
func InputGroupCheck(label string, value *hx.SignalT[bool]) hx.INode {
	return InputGroupText(hx.Input().Type("checkbox").Class("form-check-input", "mt-0").
		Attribute("aria-label", label).BindChecked(value).BindOnChecked(value))
}

// FormRow is a Row with gutters, for a form of several columns.
func FormRow(children ...hx.INode) hx.INode { return Row(children...).Class(Gutter(3)) }
