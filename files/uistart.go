package files

import (
	ui "github.com/gizak/termui"
)

func UIStart(filename string) {
	ui.Clear()

	netInfoMap := GetInfoMapMap(filename)

	parsMap := make(map[string]*ui.Par, len(netInfoMap))
	for k, v := range netInfoMap {
		parse, lines := getPar(v)
		par := ui.NewPar(parse)
		par.Height = lines + 2
		par.BorderLabel = k

		parsMap[k] = par
	}

	// build layout
	parTime := ui.NewPar(getTime())
	parTime.Height = 1
	parTime.Border = false
	ui.Body.AddRows(
		ui.NewRow(
			ui.NewCol(12, 0, parTime)))

	for _, par := range parsMap {
		ui.Body.AddRows(
			ui.NewRow(
				ui.NewCol(12, 0, par)))
	}

	// calculate layout
	ui.Body.Align()
	ui.Render(ui.Body)

	ui.Handle("/timer/1s", func(e ui.Event) {
		t := e.Data.(ui.EvtTimer)
		i := t.Count
		if i > 103 {
			ui.StopLoop()
			return
		}
		netInfoMap := GetInfoMapMap(filename)
		ui.Clear()

		parTime.Text = getTime()
		for k, v := range netInfoMap {
			parse, lines := getPar(v)

			parsMap[k].Text = parse
			parsMap[k].Height = lines + 2
		}

		ui.Render(ui.Body)

	})
}
