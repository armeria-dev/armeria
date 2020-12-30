package armeria

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
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

type TextOperation struct {
	Text string
}

// TextStyle will style text according to one or more styling options.
func TextStyle(text interface{}, opts ...TextOperation) string {
	t := fmt.Sprintf("%v", text)

	for _, o := range opts {
		t = fmt.Sprintf(o.Text, t)
	}

	return t
}

// WithBold formats the text as bold.
func WithBold() TextOperation {
	return TextOperation{
		Text: "<span style='font-weight:600'>%v</span>",
	}
}

// WithItalics formats the text using italics.
func WithItalics() TextOperation {
	return TextOperation{
		Text: "<span style='font-style:italic'>%v</span>",
	}
}

// WithMonospace formats the text using a monospace font.
func WithMonospace() TextOperation {
	return TextOperation{
		Text: "<span class='monospace'>%v</span>",
	}
}

// WithButton formats the text creating a clickable button (with optional promptData).
func WithButton(cmd, promptData string) TextOperation {
	return TextOperation{
		Text: "<span class='inline-button' data-cmd='" + cmd + "' data-prompt='" + promptData + "'>%v</span>",
	}
}

// WithLinkCmd formats the text creating a hyperlink that executes a specific command when clicked on.
func WithLinkCmd(cmd string) TextOperation {
	enc := base64.StdEncoding.EncodeToString([]byte(cmd))
	return TextOperation{
		Text: "<a href='#' class='inline-command' " +
			"data-command='" + enc + "' " +
			"tooltip='Run: " + cmd + "'>" +
			"%v</a>",
	}
}

// WithColor formats the text using a specific color.
func WithColor(color string) TextOperation {
	return TextOperation{
		Text: "<span style='color:#" + color + "'>%v</span>",
	}
}

// WithColor formats the text using a specific color.
func WithUserColor(c *Character, color int) TextOperation {
	return TextOperation{
		Text: "<span style='color:" + c.UserColor(color) + "'>%v</span>",
	}
}

// WithLink formats the text creating a hyperlink.
func WithLink(url string) TextOperation {
	return TextOperation{
		Text: "<a href='" + url + "' class='inline-link' target='_new'>%v</a>",
	}
}

// WithSize formats the text using a specific size.
func WithSize(size int) TextOperation {
	return TextOperation{
		Text: "<span style='font-size:" + strconv.Itoa(size) + "px'>%v</span>",
	}
}

// WithItemTooltip formats the text allowing a player to mouse-over the item and view the item tooltip.
func WithItemTooltip(uuid string) TextOperation {
	return TextOperation{
		Text: "<span class='hover-item-tooltip' data-uuid='" + uuid + "'>%v</span>",
	}
}

// WithContextMenu formats the text to display a context menu when right-clicking.
func WithContextMenu(name, objType, color string, content []string) TextOperation {
	enc := base64.StdEncoding.EncodeToString([]byte(strings.Join(content, ";")))
	return TextOperation{
		Text: fmt.Sprintf(
			"<span class='dynamic-context-menu' data-name='%s' data-type='%s' data-color='%s' data-content='%s'>%%v</span>",
			name,
			objType,
			color,
			enc,
		),
	}
}

// WithConvoSelection formats the text as a conversation answer.
func WithConvoSelection(id, mobUUID string, groupId int64) TextOperation {
	return TextOperation{
		Text: fmt.Sprintf(
			"<span class='convo-select' data-group='%d' data-convo-option-id='%s' data-mob-uuid='%s'>%%v</span>",
			groupId,
			id,
			mobUUID,
		),
	}
}

// WithChannelLabel formats the text as a channel header label.
func WithChannelLabel(color string) TextOperation {
	return TextOperation{
		Text: fmt.Sprintf(
			"<span style='background-color:%s;padding:3px 8px;color:#fff;border-radius:12px;margin-right:2px'>%%v</span>",
			color,
		),
	}
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

// TextCapitalization will return the same string with the first character always capitalized.
func TextCapitalization(text string) string {
	return strings.ToUpper(text[0:1]) + text[1:]
}

// TextTable returns a table in HTML with rows generated by TableRow.
func TextTable(rows ...string) string {
	rowString := ""
	for _, row := range rows {
		rowString = rowString + row
	}
	return "<table cellspacing=\"0\">" + rowString + "</table"
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
