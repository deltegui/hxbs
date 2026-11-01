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
