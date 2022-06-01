package utils

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func CrateNormalCell(content string) *tview.TableCell {
	return tview.NewTableCell(content).SetAlign(tview.AlignCenter)
}

func CrateCell(content string, color tcell.Color) *tview.TableCell {
	return tview.NewTableCell(content).SetTextColor(color).SetAlign(tview.AlignCenter)
}
