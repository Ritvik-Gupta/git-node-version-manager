package tui

import (
	"fmt"

	"github.com/Ritvik-Gupta/git-node-version-manager/utils"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func (data *TuiApplication) addPackagePages() {
	numPages := data.pages.GetPageCount()

	for pkgIdx, pkgReq := range data.packagesRequired {
		table := tview.NewTable()
		table.SetBorder(true)
		table.SetBorders(true)
		table.SetBorderPadding(0, 0, 20, 20)

		table.SetTitle(fmt.Sprintf("%s @ %s", pkgReq.Name, pkgReq.GetVersion()))
		table.SetCell(0, 0, utils.CrateCell("Repo Name", tcell.ColorYellow))
		table.SetCell(0, 1, utils.CrateCell("Version Found", tcell.ColorYellow))
		table.SetCell(0, 2, utils.CrateCell("Does Satisfy ?", tcell.ColorYellow))
		table.SetCell(0, 3, utils.CrateCell("Reason", tcell.ColorYellow))

		repoIdx := 1
		for _, repo := range data.repositories {
			table.SetCell(repoIdx, 0, utils.CrateNormalCell(repo.Name))
			table.SetCell(repoIdx, 1, utils.CrateCell(repo.Packages[pkgIdx].GetVersion(), tcell.ColorYellow))
			doesSatisfyVersion := "no"
			color := tcell.ColorRed
			versionComparision := repo.Packages[pkgIdx].VersionRankComparedTo(&pkgReq)

			if -1 <= versionComparision && versionComparision <= 2 {
				doesSatisfyVersion = "yes"
				color = tcell.ColorGreen
			}
			table.SetCell(repoIdx, 2, utils.CrateCell(doesSatisfyVersion, color))

			reason := ""
			switch versionComparision {
			case -4:
				reason = "Does not contain the Package"
			case -3:
				reason = "Is a Major Version Behind"
			case -2:
				reason = "Has a small feature missing"
			case -1:
				reason = "Only behind by a minor bug fix"
			case 1:
				reason = "Only ahead by a minor bug fix"
			case 2:
				reason = "Contains a new small feature"
			case 3:
				reason = "Is a Major Version Ahead"
			}
			table.SetCell(repoIdx, 3, utils.CrateCell(reason, color))

			repoIdx++
		}

		data.pages.AddPage(
			fmt.Sprint(numPages),
			table,
			true,
			false,
		)
		numPages++
	}
}
