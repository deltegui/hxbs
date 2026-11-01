# hxbs

Every component of [Bootstrap](https://getbootstrap.com) 5.3 as a Hix node. Only the
CSS of Bootstrap is used. What Bootstrap does with JavaScript (modal, offcanvas,
dropdown, collapse, accordion, tabs, toast, carousel, navbar toggler) is a `Signal`
here, so it works with `hx/ssr`, with `hx/wasm` and when the browser hydrates what
the server wrote. It imports only `hx`.

```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css">
```

(`bs.CSS` has that address; `bs.Icons` is Bootstrap Icons, which is a separate package.)

```go
import "github.com/deltegui/hxbs"

open := hx.Signal(false)
bs.Container(
	bs.Heading(1, "Hello"),
	bs.Button("Open", bs.Primary, bs.Lg).OnClick(func(hx.EventContext) { open.Set(true) }),
	bs.Modal(open, "Title", func() { open.Set(false) },
		func() hx.INode { return hx.P().Text("A modal is in the tree only while it is open.") },
		nil),
)
```

A component takes what it needs (text, children, the signal it reads) and returns an
`hx.INode`. A variant is a `Mod`, and the component puts its own prefix in front:
`bs.Primary` is `btn-primary` on a `Button`, `alert-primary` on an `Alert` and
`text-bg-primary` on a `Badge`. `bs.Outline(bs.Danger)` is `btn-outline-danger`.
Everything else (`mb-3`, `d-flex`, `justify-content-center`) is a class: `.Class("mb-3")`.
Helpers for the common utilities: `Space("m", "t", 3)`, `Gap`, `TextColor`, `BgColor`,
`BorderColor`, `TextBg`, `ColOn(bs.MD, 6, ...)`, `Offset`.

| File | Components |
| --- | --- |
| `layout.go` | Container(On), Row, RowCols, Col, ColSpan, ColOn, Flex, FlexColumn, Center, Spacer, Stack, Inline, Grid, HStack, VStack, Ratio, Heading, Display, Lead, Blockquote, ListUnstyled, ListInline, StretchedLink, Image, Figure, Table, SimpleTable, TableRow, TableCaption, TableResponsive |
| `form.go` | Label, FormText, FormField, FormRow, Input, TextArea, Select, FileInput, ColorInput, Datalist, PlainText, Range, Checkbox, CheckInline, CheckboxGroup, Switch, Radio(s), RadiosInline, ButtonCheck, ButtonRadios, Fieldset, ValidFeedback, IsValid, FloatingLabel, InputGroup(Text, Size, Check) |
| `components.go` | Button, SubmitButton, LoadingButton, ButtonLink, CloseButton, ButtonGroup, ButtonToolbar, Alert (+ AlertHeading, AlertLink), Badge, WithBadge, Breadcrumb, Card (+ parts, CardGroup, CardGrid, CardImageOverlay, CardHeaderTabs), ListGroup, Pagination, Progress(Of), ProgressStacked, ProgressSegment(Of), Spinner, Placeholder |
| `overlay.go` | Modal, ModalCentered, Offcanvas, Dropdown, SplitDropdown, NavDropdown, Collapse, Accordion, Toast, ToastContainer |
| `nav.go` | Tabs, TabPanel, Nav, Navbar (+ Brand, Toggler, Collapse, Nav, Link, Text), Carousel |

The layout utilities are a few helpers, not a wrapper of every class. `bs.Flex(a, bs.Spacer(), b)`
is `d-flex`; the rest is classes, with helpers that check the name:
`.Class(bs.Gap(2), bs.JustifyContent(bs.AlignBetween), bs.AlignItems(bs.AlignCenter), bs.FlexWrap)`,
`row.Class(bs.Gutter(3))`, `nav.Class(bs.FixedTop)`. `bs.Grid(3, 2, ...)` is a CSS grid of three
columns (Bootstrap has no utility for the columns of `d-grid`).

More mods: `bs.Dropup`, `bs.Dropend`, `bs.Dropstart` (where a dropdown opens) and `bs.MenuEnd`
(its menu aligns to the right) go with the variant of the button in `Dropdown`, `SplitDropdown`
and `NavDropdown`; `bs.Scrollable` and `bs.Fullscreen` go in a `Modal`.

Form notes: `ButtonCheck` and `ButtonRadios` are inputs with `btn-check` and a `label.btn`, so they are
`btn-outline-primary` unless you give a variant. Mark a control with `IsValid(control, signal)` and put a
`ValidFeedback` after it (Bootstrap shows it only when it follows the control). A button or select in an
`InputGroup` needs no wrapper; `InputGroupCheck` is the checkbox that goes in a `span.input-group-text`.

Things that differ from Bootstrap with JavaScript:

- A closed Modal, Offcanvas or Toast is not in the tree; its content is built each time it opens.
  Collapse, Accordion and Carousel build their content once and hide it with CSS.
- Controls need an `id` (`Input("text", "name", sig)`) because Bootstrap links labels by `for`.
  It must be unique in the page. `FormField` adds `is-invalid` to the control while its errors are not empty.
- Dropdown closes with a click outside (an invisible layer), but not on its own when an item is chosen.
  Its menu has `data-bs-popper="static"`, which is what makes Bootstrap's CSS place it (dropup, menu-end...) without Popper.
- `LoadingButton` shows the spinner and is disabled while its signal is true; it does not set the signal.
- In a `Table`, rows go in a `TBody` (the browser adds one if they are not, and hydration then sees a different tree).
- The modal does not lock the scroll of the page, the carousel does not slide, and there are no
  tooltips or popovers (they need Popper). Close on Escape is not handled.
- `Body` replaces the children; it does not append.

`go test ./...` (with `GOWORK=off` if there is a `go.work` above).
