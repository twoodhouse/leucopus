package memory

import (
	"math"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/truthTable"
)

type Memory struct {
	River             map[*info.Info][]int
	Paths             map[*info.Info]*pather.Path
	SupportingInfos   map[*info.Info][]*info.Info
	InfoFitGoodnesses map[*info.Info]float32 //TODO: Maybe add another lookup? both info and path are needed in case it needs to go back to an older more successful path
	ExitILinkInputs   map[*info.Info]map[*truthTable.Link]int
	Cascades          map[*info.Info][]map[*info.Info][]int //this second to last array is the list of cascades to run
	Depths            map[*info.Info]int
	RiverBalance      map[*info.Info][]int
	defaultDepth      int
}

func New() *Memory {
	var entity = Memory{
		make(map[*info.Info][]int),
		make(map[*info.Info]*pather.Path),
		make(map[*info.Info][]*info.Info),
		make(map[*info.Info]float32),
		make(map[*info.Info]map[*truthTable.Link]int),
		make(map[*info.Info][]map[*info.Info][]int),
		make(map[*info.Info]int),
		make(map[*info.Info][]int),
		10,
	}
	return &entity
}

func (m *Memory) ProcessNextIteration(values map[*info.Info]int) {
	if len(values) != len(m.River) {
		println("River row to add does not contain the same infos that the river contains")
	}
	for nfo, val := range values {
		m.AddToRiver(nfo, val)
	}
}

func (m *Memory) PrintPaths() {
	for nfo, path := range m.Paths {
		print("Path for ")
		println(nfo.Uid)
		path.Print()
	}
	println()
}

func (m *Memory) PrintGoodness() {
	println("PRINTING MEMORY FIT GOODNESS")
	for nfo, goodness := range m.InfoFitGoodnesses {
		print(nfo.Uid)
		print(": ")
		println(goodness)
	}
	println()
}

func (m *Memory) PrintNumCascades() {
	println("PRINTING NUMBER OF CASCADES")
	for nfo, cascades := range m.Cascades {
		print(nfo.Uid)
		print(": ")
		println(len(cascades))
	}
	println()
}

func (m *Memory) AddToRiver(nfo *info.Info, val int) {
	// if memory does not contain this info, add a depth for it
	if _, ok := m.Depths[nfo]; !ok {
		m.Depths[nfo] = m.defaultDepth
	}
	m.River[nfo] = append(m.River[nfo], val)
	if len(m.River[nfo]) > m.Depths[nfo] {
		m.River[nfo] = append(m.River[nfo][:0], m.River[nfo][0+1:]...) //remove first element
	}
}

func (m *Memory) GetRiverTop() map[*info.Info]int {
	riverTop := make(map[*info.Info]int)
	for info, history := range m.River {
		riverTop[info] = history[len(history)-1]
	}
	return riverTop
}

func (m *Memory) SetRiver(nfo *info.Info, vals []int) {
	// if memory does not contain this info, add a depth for it
	if _, ok := m.Depths[nfo]; !ok {
		m.Depths[nfo] = m.defaultDepth
	}
	m.River[nfo] = vals
}

func (m *Memory) GetNumCascades(nfo *info.Info) int {
	return len(m.Cascades[nfo])
}

func (m *Memory) OpenCascade(nfo *info.Info) {
	// supportingInfos := m.SupportingInfos[nfo]
	// supportingInfosUid := utility.GetUidFromInfos(supportingInfos)
	if _, ok := m.Cascades[nfo]; !ok {
		m.Cascades[nfo] = make([]map[*info.Info][]int, 0)
	}
	if _, ok := m.Cascades[nfo]; !ok {
		m.Cascades[nfo] = make([]map[*info.Info][]int, 0)
	}
	cascadeToAdd := make(map[*info.Info][]int)
	for nfo, riverRow := range m.River {
		for _, e := range riverRow {
			cascadeToAdd[nfo] = append(cascadeToAdd[nfo], e)
		}
	}
	m.Cascades[nfo] = append(m.Cascades[nfo], cascadeToAdd)
}

// //TODO: test this function
// func (m *Memory) GenerateAllCascadesTestsGivenInfoAndSupport(nfo *info.Info, supportingInfos []*info.Info, depth int) [][]map[*info.Info][]int {
// 	allTests := make([][]map[*info.Info][]int, 0)
// 	nfoCascades := m.Cascades[nfo]
// 	for _, nfoCascade := range nfoCascades {
// 		hasAllSupport := true
// 		for _, supportingInfo := range supportingInfos {
// 			if _, ok := nfoCascade[supportingInfo]; !ok {
// 				hasAllSupport = false
// 			}
// 		}
// 		if hasAllSupport {
// 			result := m.GenerateCascadeTests(nfoCascade, nfo, supportingInfos, depth)
// 			allTests = append(allTests, result)
// 		}
// 	}
// 	return allTests
// }

/*
This function needs to generate tests from a cascade, rather than river
*/
func (m *Memory) GenerateCascadeTests(cascade map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, depth int) []map[*info.Info][]int {
	tests := make([]map[*info.Info][]int, 0)
	riverDry := false
	shortestLength := len(cascade[nfo])
	for _, supportingInfo := range supportingInfos {
		if len(cascade[supportingInfo]) < shortestLength {
			shortestLength = len(cascade[supportingInfo])
		}
	}
	for i := shortestLength; !riverDry; i-- { //TODO: fix this section to take from end of list to the beginning
		riverTest := make(map[*info.Info][]int, 0)
		offset := len(cascade[nfo]) - shortestLength
		riverTest[nfo] = cascade[nfo][i-depth-1+offset : i+offset]
		for _, supportingInfo := range supportingInfos {
			offset := len(cascade[supportingInfo]) - shortestLength
			riverTest[supportingInfo] = cascade[supportingInfo][i-depth-1+offset : i+offset]
			if i-depth-1 <= 0 {
				riverDry = true
			}
		}
		// m.PrintRiverTool(riverTest)
		tests = append(tests, riverTest)
		riverTest = nil
	}

	return tests
}

func (m *Memory) ProcessPathAgainstCascades(focusInfo *info.Info, supportingInfos []*info.Info, pth *pather.Path) float32 {
	totalCount := 0
	correctCount := 0
	// Process all cascades and use correctness to gauge fitgoodness
	cascadesToProcess := m.Cascades[focusInfo]
	localRiverBalance := []int{}
	for i := 0; i < len(pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs); i++ {
		localRiverBalance = append(localRiverBalance, 0)
	}
	numCorrect := pather.ProcessAllCascadesTogether(cascadesToProcess, focusInfo, supportingInfos, pth, localRiverBalance)
	correctCount = correctCount + numCorrect
	totalCount = totalCount + len(cascadesToProcess)
	// Update fitgoodness
	var goodness float32
	if len(cascadesToProcess) == 0 {
		goodness = 1.0
	} else {
		goodness = float32(correctCount) / float32(totalCount)
	}

	return goodness

	// for _, cascade := range cascadesToProcess {
	// 	pth.TakeSnapshot()
	// 	cascadeSuccess := pather.ProcessCascadeWithIVariation(cascade, focusInfo, supportingInfos, pth)
	// 	pth.RestoreSnapshot()
	// 	totalCount = totalCount + 1
	// 	if cascadeSuccess {
	// 		correctCount = correctCount + 1
	// 	}
	// 	_ = cascade
	// }
	//
	// // Update fitgoodness
	// var unweightedGoodness float32
	// if len(cascadesToProcess) == 0 {
	// 	unweightedGoodness = 1.0
	// } else {
	// 	unweightedGoodness = float32(correctCount) / float32(totalCount)
	// }
	// // oldFitGoodness := m.InfoFitGoodnesses[focusInfo]
	// // m.InfoFitGoodnesses[focusInfo] = (7 * oldFitGoodness / 8) + (goodness / 8)
	// // println(m.InfoFitGoodnesses[focusInfo])
	// return unweightedGoodness
}

func (m *Memory) MagicRiverInput(row map[*info.Info]int) {
	for nfo, _ := range row {
		if _, ok := m.RiverBalance[nfo]; !ok {
			m.RiverBalance[nfo] = make([]int, 0)
			if _, ok := m.Paths[nfo]; ok {
				for i := 0; i < int(math.Exp2(float64(len(m.Paths[nfo].MiddleLinks[len(m.Paths[nfo].MiddleLinks)-1].Inputs)))); i++ {
					m.RiverBalance[nfo] = append(m.RiverBalance[nfo], 0)
				}
			}
		}
	}
	// if this is the first time through, just process the row and add it to the river
	for k, _ := range row {
		if len(m.River[k]) == 0 {
			m.ProcessNextIteration(row)
			return
		}
	}
	// Now ProcessRiver (for each river info) with river top
	for focusInfo, _ := range m.River {
		if _, ok := m.Paths[focusInfo]; !ok {
			continue // don't run the rest of this path related code if there is no path yet for the info
		}
		supportingInfos := m.SupportingInfos[focusInfo]
		pth := m.Paths[focusInfo]
		var result bool
		var index int
		// If the first time, it needs to set the ExitILinkInputs to 1. The path switcher in test master can do consideration on what things to set when switching
		exitILinkInputs := make(map[*truthTable.Link]int)
		if len(pth.ExitILinks) > 0 && pth.ExitILinks[0].Output == -1 {
			for _, lnk := range pth.ExitILinks {
				exitILinkInputs[lnk] = 1
			}
			result, index = pather.ProcessRiver2(m.GetRiverTop(), row, exitILinkInputs, focusInfo, supportingInfos, pth, true)
		} else {
			result, index = pather.ProcessRiver2(m.GetRiverTop(), row, exitILinkInputs, focusInfo, supportingInfos, pth, false)
		}
		// apply river balance (helps decide if there is a better explanation)
		if result {
			m.RiverBalance[focusInfo][index] = m.RiverBalance[focusInfo][index] + 1
		} else {
			m.RiverBalance[focusInfo][index] = m.RiverBalance[focusInfo][index] - 1
		}
		// check river balance to see if the explanation should be shifted. Note that this will NOT switch around what cascade must be checked
		if m.RiverBalance[focusInfo][index] < 0 {
			valToChangeTo := 0
			if pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs[index] == 0 {
				valToChangeTo = 1
			}
			pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs[index] = valToChangeTo
			m.RiverBalance[focusInfo][index] = m.RiverBalance[focusInfo][index] * -1
		}

		// Open cascade upon failure of process river
		if !result {
			// println("failure")
			// pth.Print()
			m.OpenCascade(focusInfo)
		}
		totalCount := 1
		correctCount := 0
		if result {
			correctCount = 1
		}
		// Process all cascades and use correctness to gauge fitgoodness
		cascadesToProcess := m.Cascades[focusInfo]
		numCorrect := pather.ProcessAllCascadesTogether(cascadesToProcess, focusInfo, supportingInfos, pth, m.RiverBalance[focusInfo])
		correctCount = correctCount + numCorrect
		totalCount = totalCount + len(cascadesToProcess)
		// for _, cascade := range cascadesToProcess {
		// 	pth.TakeSnapshot()
		// 	cascadeSuccess := pather.ProcessCascadeWithIVariation(cascade, focusInfo, supportingInfos, pth)
		// 	pth.RestoreSnapshot()
		// 	totalCount = totalCount + 1
		// 	if cascadeSuccess {
		// 		correctCount = correctCount + 1
		// 	}
		// 	_ = cascade
		// }

		// Update fitgoodness
		var goodness float32
		if len(cascadesToProcess) == 0 {
			goodness = 1.0
		} else {
			goodness = float32(correctCount) / float32(totalCount)
		}
		_ = goodness
		// oldFitGoodness := m.InfoFitGoodnesses[focusInfo]
		m.InfoFitGoodnesses[focusInfo] = goodness
		// println(m.InfoFitGoodnesses[focusInfo])
	}
	m.ProcessNextIteration(row)

}

func (m *Memory) PrintRiver() {
	max := 0
	for _, v := range m.Depths {
		if v > max {
			max = v
		}
	}
	for info, history := range m.River {
		for i := 0; i < max-len(m.River[info]); i++ {
			print(" ")
		}
		// reversedHistory := make([]int, len(history))
		// copy(reversedHistory, history)
		// reversedHistory = reverseInts(reversedHistory)
		// for _, val := range reversedHistory {
		for _, val := range history {
			print(val)
		}
		print(": ")
		print(info.Uid)
		println("")
	}
	println("")
}

func (m *Memory) PrintRiverBalance() {
	m.PrintRiverTool(m.RiverBalance)
}

func (m *Memory) PrintRiverTool(river map[*info.Info][]int) {
	max := 0
	for _, v := range m.Depths {
		if v > max {
			max = v
		}
	}
	for nfo, history := range river {
		for i := 0; i < max-len(river[nfo]); i++ {
			print(" ")
		}
		// reversedHistory := make([]int, len(history))
		// copy(reversedHistory, history)
		// reversedHistory = reverseInts(reversedHistory)
		// for _, val := range reversedHistory {
		for _, val := range history {
			print(val)
		}
		print(": ")
		print(nfo.Uid)
		println("")
	}
	println("")
}

func reverseInts(input []int) []int {
	if len(input) == 0 {
		return input
	}
	return append(reverseInts(input[1:]), input[0])
}

func GeneralPrintRiverTool(river map[*info.Info][]int) {
	max := 20
	for nfo, history := range river {
		for i := 0; i < max-len(river[nfo]); i++ {
			print(" ")
		}
		// reversedHistory := make([]int, len(history))
		// copy(reversedHistory, history)
		// reversedHistory = reverseInts(reversedHistory)
		// for _, val := range reversedHistory {
		for _, val := range history {
			print(val)
		}
		print(": ")
		print(nfo.Uid)
		println("")
	}
	println("")
}

// func TestPath(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, links []*truthTable.Link) bool {
// 	//TODO: create pathing model in order to have this fed into Test also.
// 	return false
// }
//
// func TestPathHelper(inputs []int, table *truthTable.TruthTable) int {
// 	return 0
// }
