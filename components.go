package bs

import (
	"strconv"

	"github.com/deltegui/hx"
)

// Button is a button that does not submit. Without mods it is Primary. Mods:
// colors, Outline(color), Sm, Lg.
func Button(label string, mods ...Mod) hx.INode { return button(label, "button", mods) }

// SubmitButton submits the form it is in.
func SubmitButton(label string, mods ...Mod) hx.INode { return button(label, "submit", mods) }

// ButtonLink is an anchor that looks like a button.
func ButtonLink(href, label string, mods ...Mod) hx.INode {
	a := hx.A().Href(href).Class("btn").Attribute("role", "button").Text(label)
	if len(mods) == 0 {
		return a.Class("btn-primary")
	}
	return with(a, "btn", mods)
}

// CloseButton is the "x" of alerts, modals and toasts.
func CloseButton(onClick func()) hx.INode {
	b := hx.Button().Type("button").Class("btn-close").Attribute("aria-label", "Close")
	if onClick != nil {
		b.OnClick(func(hx.EventContext) { onClick() })
	}
	return b
}

// ButtonGroup sticks buttons together. Mods: Sm, Lg, Vertical.
func ButtonGroup(mods []Mod, children ...hx.INode) hx.INode {
	g := hx.Div().Class("btn-group").Attribute("role", "group").Body(children...)
	for _, m := range mods {
		if m == Vertical {
			g.RemoveClass("btn-group").Class("btn-group-vertical")
		} else {
			g.Class("btn-group-" + string(m))
		}
	}
	return g
}

// ButtonToolbar puts several ButtonGroups in a line.
func ButtonToolbar(children ...hx.INode) hx.INode {
	return hx.Div().Class("btn-toolbar").Attribute("role", "toolbar").Body(children...)
}

// Alert is a colored box that stands out. Mods: colors.
func Alert(child hx.INode, mods ...Mod) hx.INode {
	return with(hx.Div().Class("alert").Attribute("role", "alert").Body(child), "alert", mods)
}

// DismissibleAlert is an Alert with a close button in the corner.
func DismissibleAlert(onClose func(), child hx.INode, mods ...Mod) hx.INode {
	a := with(hx.Div().Class("alert", "alert-dismissible").Attribute("role", "alert"), "alert", mods)
	return a.Body(child, CloseButton(onClose))
}

// Badge is a small label. Mods: colors; Pill makes it round.
func Badge(text string, mods ...Mod) hx.INode {
	b := hx.Span().Class("badge").Text(text)
	for _, m := range mods {
		if m == Pill {
			b.Class("rounded-pill")
		} else {
			b.Class("text-bg-" + string(m))
		}
	}
	return b
}

// Pill is the mod of a Badge with round ends.
const Pill Mod = "pill"

// Breadcrumb items. The last one is usually Current.
type Crumb struct {
	Label, Href string
	Current     bool
}

// Breadcrumb is the trail "Home / Products / Edit".
func Breadcrumb(crumbs ...Crumb) hx.INode {
	items := make([]hx.INode, len(crumbs))
	for i, c := range crumbs {
		li := hx.Li().Class("breadcrumb-item")
		if c.Current {
			li.Class("active").Attribute("aria-current", "page").Text(c.Label)
		} else {
			li.Body(hx.A().Href(c.Href).Text(c.Label))
		}
		items[i] = li
	}
	return hx.Nav().Attribute("aria-label", "breadcrumb").Body(hx.Ol().Class("breadcrumb").BodyList(items))
}

// Card is a box with parts: put CardHeader, CardImage, CardBody and CardFooter in
// it. Chain TextBg(Primary) or BorderColor(Danger) for a color.
func Card(children ...hx.INode) hx.INode { return hx.Div().Class("card").Body(children...) }

func CardHeader(children ...hx.INode) hx.INode {
	return hx.Div().Class("card-header").Body(children...)
}

func CardBody(children ...hx.INode) hx.INode { return hx.Div().Class("card-body").Body(children...) }

func CardFooter(children ...hx.INode) hx.INode {
	return hx.Div().Class("card-footer").Body(children...)
}

// CardImage is the image of the top of a card.
func CardImage(src, alt string) hx.INode {
	return hx.Img().Src(src).Class("card-img-top").Attribute("alt", alt)
}

func CardTitle(text string) hx.INode { return hx.H5().Class("card-title").Text(text) }

func CardSubtitle(text string) hx.INode {
	return hx.H6().Class("card-subtitle", "mb-2", "text-body-secondary").Text(text)
}

func CardText(text string) hx.INode { return hx.P().Class("card-text").Text(text) }

// ListGroup is a list of items. Mods: Flush, Numbered, Horizontal.
func ListGroup(mods []Mod, items ...hx.INode) hx.INode {
	l := hx.Ul().Class("list-group").Body(items...)
	return with(l, "list-group", mods)
}

// ListGroupItem is an item of a ListGroup. Mods: colors. Chain "active" or
// "disabled" on it.
func ListGroupItem(child hx.INode, mods ...Mod) hx.INode {
	return with(hx.Li().Class("list-group-item").Body(child), "list-group-item", mods)
}

// ListGroupAction is an item that is a button, with a hover.
func ListGroupAction(label string, onClick func(), mods ...Mod) hx.INode {
	b := hx.Button().Type("button").Class("list-group-item", "list-group-item-action").Text(label)
	if onClick != nil {
		b.OnClick(func(hx.EventContext) { onClick() })
	}
	return with(b, "list-group-item", mods)
}

// Pagination shows previous, next, and a window of pages around the current one
// with the first and the last. page and pages count from 1. Mods: Sm, Lg.
func Pagination(page, pages hx.Gettable[int], goTo func(int), mods ...Mod) hx.INode {
	item := func(label string, target int, active, disabled bool) hx.INode {
		li := hx.Li().Class("page-item")
		b := hx.Button().Type("button").Class("page-link").Text(label)
		switch {
		case active:
			li.Class("active")
			b.Attribute("aria-current", "page")
		case disabled:
			li.Class("disabled")
			b.Attribute("disabled", "")
		default:
			b.OnClick(func(hx.EventContext) { goTo(target) })
		}
		return li.Body(b)
	}
	gap := func() hx.INode {
		return hx.Li().Class("page-item", "disabled").Body(hx.Span().Class("page-link").Text("…"))
	}
	list := hx.Ul().Class("pagination")
	with(list, "pagination", mods)
	return hx.Nav().Attribute("aria-label", "pagination").Body(list.Body(hx.Dynamic(func() hx.INode {
		current, last := page.Get(), pages.Get()
		items := []hx.INode{item("«", current-1, false, current <= 1)}
		from, to := max(1, current-2), min(last, current+2)
		if from > 1 {
			items = append(items, item("1", 1, current == 1, false))
			if from > 2 {
				items = append(items, gap())
			}
		}
		for p := from; p <= to; p++ {
			items = append(items, item(strconv.Itoa(p), p, p == current, false))
		}
		if to < last {
			if to < last-1 {
				items = append(items, gap())
			}
			items = append(items, item(strconv.Itoa(last), last, current == last, false))
		}
		items = append(items, item("»", current+1, false, current >= last))
		return hx.Noop().BodyList(items)
	})))
}

// Progress is a bar that is value of max full. Mods on the bar: colors,
// Striped, Animated.
func Progress(value, max int, mods ...Mod) hx.INode {
	bar := progressBar(mods)
	setProgress(bar, value, max)
	return progressBox(bar, value, max)
}

// ProgressOf is a Progress that follows a signal.
func ProgressOf(value hx.Gettable[int], max int, mods ...Mod) hx.INode {
	bar := progressBar(mods)
	box := progressBox(bar, 0, max)
	return box.Effect(func() {
		v := value.Get()
		setProgress(bar, v, max)
		box.Attribute("aria-valuenow", strconv.Itoa(v))
	})
}

func progressBar(mods []Mod) hx.INode {
	bar := hx.Div().Class("progress-bar")
	for _, m := range mods {
		switch m {
		case Striped, Animated:
			bar.Class("progress-bar-" + string(m))
			if m == Animated {
				bar.Class("progress-bar-striped")
			}
		default:
			bar.Class("text-bg-" + string(m))
		}
	}
	return bar
}

func setProgress(bar hx.INode, value, max int) {
	pct := 0
	if max > 0 {
		pct = value * 100 / max
	}
	bar.Style("width", strconv.Itoa(pct)+"%")
}

func progressBox(bar hx.INode, value, max int) hx.INode {
	return hx.Div().Class("progress").Attribute("role", "progressbar").
		Attribute("aria-valuenow", strconv.Itoa(value)).Attribute("aria-valuemin", "0").
		Attribute("aria-valuemax", strconv.Itoa(max)).Body(bar)
}

// Spinner is the circle that turns while something loads. Mods: colors (their
// text color), Sm, Grow (a pulsing dot instead of a ring).
func Spinner(mods ...Mod) hx.INode {
	kind := "spinner-border"
	for _, m := range mods {
		if m == Grow {
			kind = "spinner-grow"
		}
	}
	s := hx.Div().Class(kind).Attribute("role", "status").Body(hx.Span().Class("visually-hidden").Text("Loading..."))
	for _, m := range mods {
		switch m {
		case Grow:
		case Sm:
			s.Class(kind + "-sm")
		default:
			s.Class("text-" + string(m))
		}
	}
	return s
}

// Placeholder is a gray bar that stands in for text that loads. width is the
// number of columns of 12 it takes.
func Placeholder(width int) hx.INode {
	return hx.Span().Class("placeholder", "col-"+strconv.Itoa(width))
}

// PlaceholderGlow wraps Placeholders to make them pulse.
func PlaceholderGlow(children ...hx.INode) hx.INode {
	return hx.P().Class("placeholder-glow").Body(children...)
}

// AlertHeading and AlertLink are the title and the links inside an Alert.
func AlertHeading(text string) hx.INode { return hx.H4().Class("alert-heading").Text(text) }
func AlertLink(href, text string) hx.INode {
	return hx.A().Href(href).Class("alert-link").Text(text)
}

// CardGroup joins cards in a row with no space between them.
func CardGroup(children ...hx.INode) hx.INode { return hx.Div().Class("card-group").Body(children...) }

// CardGrid puts each card in a column: one column on a phone, cols from md up.
func CardGrid(cols int, cards ...hx.INode) hx.INode {
	cells := make([]hx.INode, len(cards))
	for i, c := range cards {
		cells[i] = hx.Div().Class("col").Body(c)
	}
	return hx.Div().Class("row", "row-cols-1", "row-cols-md-"+strconv.Itoa(cols), "g-4").BodyList(cells)
}

// CardImageOverlay is the text that goes over a card image.
func CardImageOverlay(children ...hx.INode) hx.INode {
	return hx.Div().Class("card-img-overlay").Body(children...)
}

// CardHeaderTabs is a CardHeader with tabs (or bs.Pills) in it.
func CardHeaderTabs(active *hx.SignalT[string], tabs []Tab, mods ...Mod) hx.INode {
	class := "card-header-tabs"
	for _, m := range mods {
		if m == Pills {
			class = "card-header-pills"
		}
	}
	return CardHeader(Tabs(active, tabs, mods...).Class(class))
}

// WithBadge puts a counter badge on the corner of a node (a button, an icon).
// Without a variant the badge is danger.
func WithBadge(target hx.INode, text string, mods ...Mod) hx.INode {
	badge := hx.Span().Class("position-absolute", "top-0", "start-100", "translate-middle", "badge", "rounded-pill").Text(text)
	if len(mods) == 0 {
		badge.Class("text-bg-danger")
	} else {
		with(badge, "text-bg", mods)
	}
	return hx.Span().Class("position-relative", "d-inline-block").Body(target, badge)
}

// LoadingButton is a Button that shows a spinner and is disabled while loading is true.
func LoadingButton(label string, loading hx.Gettable[bool], mods ...Mod) hx.INode {
	b := hx.Button().Type("button").Class("btn")
	if len(mods) == 0 {
		b.Class("btn-primary")
	} else {
		with(b, "btn", mods)
	}
	b.Body(
		hx.Dynamic(func() hx.INode {
			if !loading.Get() {
				return nil
			}
			return hx.Span().Class("spinner-border", "spinner-border-sm", "me-2").
				Attribute("role", "status").Attribute("aria-hidden", "true")
		}),
		hx.Span().Text(label))
	return b.Effect(func() {
		if loading.Get() {
			b.Attribute("disabled", "")
		} else {
			b.RemoveAttribute("disabled")
		}
	})
}

// ProgressStacked puts several ProgressSegment in one bar.
func ProgressStacked(segments ...hx.INode) hx.INode {
	return hx.Div().Class("progress-stacked").Body(segments...)
}

// ProgressSegment is a part of a ProgressStacked: value out of max of the whole bar.
func ProgressSegment(value, max int, mods ...Mod) hx.INode {
	return segmentBox(progressBar(mods), value, max)
}

// ProgressSegmentOf is ProgressSegment that follows a signal.
func ProgressSegmentOf(value hx.Gettable[int], max int, mods ...Mod) hx.INode {
	box := segmentBox(progressBar(mods), 0, max)
	return box.Effect(func() {
		v := value.Get()
		box.Style("width", strconv.Itoa(percent(v, max))+"%")
		box.Attribute("aria-valuenow", strconv.Itoa(v))
	})
}

func percent(value, max int) int {
	if max <= 0 {
		return 0
	}
	return value * 100 / max
}

func segmentBox(bar hx.INode, value, max int) hx.INode {
	return hx.Div().Class("progress").Attribute("role", "progressbar").
		Attribute("aria-valuenow", strconv.Itoa(value)).Attribute("aria-valuemin", "0").
		Attribute("aria-valuemax", strconv.Itoa(max)).
		Style("width", strconv.Itoa(percent(value, max))+"%").Body(bar)
}
