package main

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

type tabler struct {
	upperLeftCorner, upperRightCorner        string
	lowerLeftCorner, lowerRightCorner        string
	vertical, verticalLeft, verticalRight    string
	horizontal, horizontalUp, horizontalDown string
	horizontalVertical                       string // aka cross
}

func NewTabler() *tabler {
	return &tabler{
		upperRightCorner:   "╗",
		lowerLeftCorner:    "╚",
		lowerRightCorner:   "╝",
		vertical:           "║",
		verticalLeft:       "╣",
		verticalRight:      "╠",
		horizontal:         "═",
		horizontalUp:       "╩",
		horizontalDown:     "╦",
		upperLeftCorner:    "╔",
		horizontalVertical: "╬",
	}
}

type textSpacingFunc func(columns int, columnWidth int) (leading, trailing int)

var (
	centerSpaceFunc = func(textWidth, columnWidth int) (int, int) {
		leading := (columnWidth - textWidth) / 2
		trailing := (columnWidth-textWidth)/2 + (columnWidth-textWidth)%2 // assume in case of inequal parity, trailing spaces would be one more

		return leading, trailing
	}

	defaultSpaceFunc = func(textWidth, columnWidth int) (int, int) {
		return 1, columnWidth - textWidth - 1
	}
)

func (t *tabler) Build(headers []string, rows [][]string) (string, error) {
	columnWidths, err := t.calculateColumnWidth(headers, rows)
	if err != nil {
		return "", errors.Wrap(err, "calculate column width")
	}

	return t.draw(headers, rows, columnWidths), nil
}

func (t *tabler) calculateColumnWidth(headers []string, rows [][]string) ([]int, error) {
	columnCount := len(headers)
	columnWidths := make([]int, columnCount)

	// Set initial column width
	for i := 0; i < columnCount; i++ {
		columnWidths[i] = len(headers[i])
	}

	// Set final column width depending on max symbol count
	for _, row := range rows {
		if len(row) != columnCount {
			return nil, errors.New(fmt.Sprintf("column count doesn't match the row count (headers: %v, row: %v)", headers, row))
		}

		for columnNr := range row {
			columnLength := len(row[columnNr])
			if columnLength > columnWidths[columnNr] {
				columnWidths[columnNr] = columnLength
			}
		}
	}

	// Include leading and trailing space
	for i := 0; i < len(columnWidths); i++ {
		columnWidths[i] += 2
	}

	return columnWidths, nil
}

func (t *tabler) draw(headers []string, rows [][]string, columnWidths []int) string {
	// Build header
	str := t.upperLeftCorner
	for idx, columnWidth := range columnWidths {
		str += strings.Repeat(t.horizontal, columnWidth)
		if idx == len(columnWidths)-1 {
			str += t.upperRightCorner + "\n"
			break
		}
		str += t.horizontalDown
	}

	str += t.assembleTableRow(headers, columnWidths, centerSpaceFunc)

	str += t.verticalRight
	for idx, columnWidth := range columnWidths {
		str += strings.Repeat(t.horizontal, columnWidth)
		if idx == len(columnWidths)-1 {
			str += t.verticalLeft + "\n"
			break
		}
		str += t.horizontalVertical
	}

	for _, row := range rows {
		str += t.assembleTableRow(row, columnWidths, defaultSpaceFunc)
	}

	str += t.lowerLeftCorner
	for idx, columnWidth := range columnWidths {
		str += strings.Repeat(t.horizontal, columnWidth)
		if idx == len(columnWidths)-1 {
			str += t.lowerRightCorner + "\n"
			break
		}
		str += t.horizontalUp
	}
	return str
}

func (t *tabler) assembleTableRow(columns []string, columnWidths []int, spacingFunc textSpacingFunc) string {
	str := t.vertical
	for idx, column := range columns {
		leadingSpaces, trailingSpaces := spacingFunc(len(column), columnWidths[idx])
		str += strings.Repeat(" ", leadingSpaces) + column + strings.Repeat(" ", trailingSpaces) + t.vertical

		if idx == len(columnWidths)-1 {
			str += "\n"
			break
		}
	}

	return str
}
