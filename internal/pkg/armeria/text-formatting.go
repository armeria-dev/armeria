package armeria

import "fmt"

const (
	TextStyleMonospace int = iota
	TextStyleBold
	TextStyleColor
	TextStyleLink
)

const (
	TextStatement int = iota
	TextQuestion
	TextExclaim
)

type TableCell struct {
	content string
	styling string
	header  bool
}

// TextStyle will style text according to one or more styling options.
func TextStyle(text interface{}, opts ...int) string {
	t := fmt.Sprintf("%v", text)

	var color string
	if t[0:1] == "#" && len(t) > 6 {
		color = t[1:7]
		t = t[7:]
	}

	for _, o := range opts {
		switch o {
		case TextStyleBold:
			t = fmt.Sprintf("<span style='font-weight:600'>%v</span>", t)
		case TextStyleMonospace:
			t = fmt.Sprintf("<span class='monospace'>%v</span>", t)
		case TextStyleColor:
			t = fmt.Sprintf("<span style='color:#%s'>%v</span>", color, t)
		case TextStyleLink:
			t = fmt.Sprintf("<a href='%[1]v' class='inline-link' target='_new'>%[1]v</a>", t)
		}
	}

	return t
}

// TextPunctuation will automatically punctuate a string and return the punctuation type.
func TextPunctuation(text string) (string, int) {
	lastChar := text[len(text)-1:]

	if lastChar == "." {
		return text, TextStatement
	} else if lastChar == "?" {
		return text, TextQuestion
	} else if lastChar == "!" {
		return text, TextExclaim
	} else {
		return text + ".", TextStatement
	}
}

// TextTable returns a table in HTML with rows generated by TableRow.
func TextTable(rows ...string) string {
	rowString := ""
	for _, row := range rows {
		rowString = rowString + row
	}
	return "<table>" + rowString + "</table"
}

// TableRow generates a row to be used within a TextTable.
func TableRow(cells ...TableCell) string {
	cellString := ""
	for _, cell := range cells {
		if cell.header {
			cellString = fmt.Sprintf("%s<th style='%s'>%s</th>", cellString, cell.styling, cell.content)
		} else {
			cellString = fmt.Sprintf("%s<td style='%s'>%s</td>", cellString, cell.styling, cell.content)
		}
	}
	return "<tr>" + cellString + "</tr>"
}
