package bs

import (
	"strconv"

	"github.com/deltegui/hx"
)

// Container is the centered box of a page. Chain "container-fluid" or
// "container-md" on it for the other widths.
func Container(children ...hx.INode) hx.INode { return hx.Div().Class("container").Body(children...) }

// FluidContainer takes the full width.
func FluidContainer(children ...hx.INode) hx.INode {
	return hx.Div().Class("container-fluid").Body(children...)
}

// Row is a row of the 12 column grid. Put Cols in it. Chain Gap(3), "g-3",
// "justify-content-center", "align-items-center"...
func Row(children ...hx.INode) hx.INode { return hx.Div().Class("row").Body(children...) }

// RowCols is a row with n columns per line: RowCols(3, ...) is "row row-cols-3".
func RowCols(n int, children ...hx.INode) hx.INode {
	return Row(children...).Class("row-cols-" + strconv.Itoa(n))
}

// Col is a column that shares the row with the others.
func Col(children ...hx.INode) hx.INode { return hx.Div().Class("col").Body(children...) }

// ColSpan is a column that takes n of the 12 parts.
func ColSpan(n int, children ...hx.INode) hx.INode {
	return hx.Div().Class("col-" + strconv.Itoa(n)).Body(children...)
}

// ColOn is a column that takes n of the 12 parts from the breakpoint on:
// ColOn(MD, 6, ...) is "col-md-6".
func ColOn(breakpoint string, n int, children ...hx.INode) hx.INode {
	return hx.Div().Class("col-" + breakpoint + "-" + strconv.Itoa(n)).Body(children...)
}

// Offset is the class that skips n parts: .Class(bs.Offset(2)).
func Offset(n int) string { return "offset-" + strconv.Itoa(n) }

// HStack and VStack put their children in a line or in a column with the gap of
// the class that you chain (Gap(3)).
func HStack(children ...hx.INode) hx.INode { return hx.Div().Class("hstack").Body(children...) }
func VStack(children ...hx.INode) hx.INode { return hx.Div().Class("vstack").Body(children...) }

// Ratio keeps the proportion of its child (an iframe, a video): "1x1", "4x3",
// "16x9" or "21x9".
func Ratio(ratio string, child hx.INode) hx.INode {
	return hx.Div().Class("ratio", "ratio-"+ratio).Body(child)
}

// Heading is a heading of level 1 to 6 (<h1>..<h6>).
func Heading(level int, text string) hx.INode {
	switch level {
	case 1:
		return hx.H1().Text(text)
	case 2:
		return hx.H2().Text(text)
	case 3:
		return hx.H3().Text(text)
	case 4:
		return hx.H4().Text(text)
	case 5:
		return hx.H5().Text(text)
	}
	return hx.H6().Text(text)
}

// Display is a big heading of level 1 to 6: it has the size of "display-1".
func Display(level int, text string) hx.INode {
	return Heading(1, text).Class("display-" + strconv.Itoa(level))
}

// Lead is a paragraph that stands out.
func Lead(text string) hx.INode { return hx.P().Class("lead").Text(text) }

// Image is an image that you can make responsive ("img-fluid") or a thumbnail
// ("img-thumbnail") by chaining the class.
func Image(src, alt string) hx.INode {
	return hx.Img().Src(src).Attribute("alt", alt)
}

// Figure is an image with a caption.
func Figure(image hx.INode, caption string) hx.INode {
	return hx.El("figure").Class("figure").Body(
		image.Class("figure-img", "img-fluid", "rounded"),
		hx.El("figcaption").Class("figure-caption").Text(caption))
}

// Table is a table. Put THead, TBody and TFoot in it, or use SimpleTable. Mods:
// Striped, Hover, Bordered, Borderless, Sm, and the colors (Dark).
func Table(mods []Mod, children ...hx.INode) hx.INode {
	return with(hx.Table().Class("table").Body(children...), "table", mods)
}

// SimpleTable is a table with a header row of headers and body inside a TBody.
func SimpleTable(mods []Mod, headers []string, rows ...hx.INode) hx.INode {
	cells := make([]hx.INode, len(headers))
	for i, h := range headers {
		cells[i] = hx.Th().Attribute("scope", "col").Text(h)
	}
	return Table(mods, hx.THead().Body(hx.Tr().BodyList(cells)), hx.TBody().Body(rows...))
}

// TableResponsive lets a wide table scroll sideways instead of breaking the page.
func TableResponsive(table hx.INode) hx.INode {
	return hx.Div().Class("table-responsive").Body(table)
}

// Align is a value of the flexbox utilities (justify-content-*, align-items-*).
type Align string

const (
	AlignStart    Align = "start"
	AlignEnd      Align = "end"
	AlignCenter   Align = "center"
	AlignBetween  Align = "between"
	AlignAround   Align = "around"
	AlignEvenly   Align = "evenly"
	AlignStretch  Align = "stretch"
	AlignBaseline Align = "baseline"
)

// FixedTop, FixedBottom and StickyTop are the position utilities of a navbar:
// nav.Class(bs.FixedTop).
const (
	FixedTop    = "fixed-top"
	FixedBottom = "fixed-bottom"
	StickyTop   = "sticky-top"
	FlexWrap    = "flex-wrap"
)

// JustifyContent is the class justify-content-*: nav.Class(bs.JustifyContent(bs.AlignBetween)).
func JustifyContent(a Align) string { return "justify-content-" + string(a) }

// AlignItems is the class align-items-*.
func AlignItems(a Align) string { return "align-items-" + string(a) }

// Gutter, GutterX and GutterY are the gutters of a Row: row.Class(bs.Gutter(3)).
func Gutter(n int) string  { return "g-" + strconv.Itoa(n) }
func GutterX(n int) string { return "gx-" + strconv.Itoa(n) }
func GutterY(n int) string { return "gy-" + strconv.Itoa(n) }

// ContainerOn is a container that is 100% wide until the breakpoint (container-md).
func ContainerOn(breakpoint string, children ...hx.INode) hx.INode {
	return hx.Div().Class("container-" + breakpoint).Body(children...)
}

// Flex is d-flex. Add the rest with classes: .Class(bs.Gap(2), bs.AlignItems(bs.AlignCenter), bs.FlexWrap).
func Flex(children ...hx.INode) hx.INode { return hx.Div().Class("d-flex").Body(children...) }

// FlexColumn is d-flex flex-column.
func FlexColumn(children ...hx.INode) hx.INode {
	return hx.Div().Class("d-flex", "flex-column").Body(children...)
}

// Center puts its children in the middle, in both directions.
func Center(children ...hx.INode) hx.INode {
	return Flex(children...).Class(JustifyContent(AlignCenter), AlignItems(AlignCenter))
}

// Spacer takes the free space of a flex container, pushing what follows it to the end.
func Spacer() hx.INode { return hx.Div().Class("flex-grow-1") }

// Stack is a vertical stack with a gap (vstack gap-N); Inline is the horizontal one.
func Stack(gap int, children ...hx.INode) hx.INode  { return VStack(children...).Class(Gap(gap)) }
func Inline(gap int, children ...hx.INode) hx.INode { return HStack(children...).Class(Gap(gap)) }

// Grid is a CSS grid of cols equal columns (d-grid with grid-template-columns, which
// the Bootstrap CSS does not have as a utility).
func Grid(cols, gap int, children ...hx.INode) hx.INode {
	return hx.Div().Class("d-grid", Gap(gap)).
		Style("grid-template-columns", "repeat("+strconv.Itoa(cols)+", 1fr)").Body(children...)
}

// Blockquote is a quotation with its source (figure, blockquote and figcaption).
func Blockquote(quote, source string) hx.INode {
	quoted := hx.Blockquote().Class("blockquote").Body(hx.P().Text(quote))
	if source == "" {
		return hx.Figure().Body(quoted)
	}
	return hx.Figure().Body(quoted, hx.Figcaption().Class("blockquote-footer").Text(source))
}

func items(class, itemClass string, nodes []hx.INode) hx.INode {
	lis := make([]hx.INode, len(nodes))
	for i, n := range nodes {
		li := hx.Li().Body(n)
		if itemClass != "" {
			li.Class(itemClass)
		}
		lis[i] = li
	}
	return hx.Ul().Class(class).BodyList(lis)
}

// ListUnstyled is a list with no bullets; each child goes in its own <li>.
func ListUnstyled(children ...hx.INode) hx.INode { return items("list-unstyled", "", children) }

// ListInline is a list in one line; each child goes in its own <li>.
func ListInline(children ...hx.INode) hx.INode {
	return items("list-inline", "list-inline-item", children)
}

// StretchedLink is a link whose click area is the whole positioned ancestor (a card).
func StretchedLink(href, label string) hx.INode {
	return hx.A().Href(href).Class("stretched-link").Text(label)
}

// TableRow is a <tr> with a variant (table-primary, table-active...).
func TableRow(mods []Mod, children ...hx.INode) hx.INode {
	return with(hx.Tr().Body(children...), "table", mods)
}

// TableCaption is the <caption> of a table.
func TableCaption(text string) hx.INode { return hx.Caption().Text(text) }
