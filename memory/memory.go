package memory

import (
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

func (m *Memory) ProcessPathAgainstCascades(pth *pather.Path) float32 {
	//TODO: write this function
	return float32(0)
}

func (m *Memory) MagicRiverInput(row map[*info.Info]int) {
	// First ProcessNextIteration(bla)
	m.ProcessNextIteration(row)
	// Now ProcessRiver (for each river info) with river top
	for focusInfo, _ := range m.River {
		supportingInfos := m.SupportingInfos[focusInfo]
		pth := m.Paths[focusInfo]
		var result bool
		// If the first time, it needs to set the ExitILinkInputs to 1. The path switcher in test master can do consideration on what things to set when switching
		exitILinkInputs := make(map[*truthTable.Link]int)
		if m.Paths[focusInfo].ExitILinks[0].Output == -1 {
			for _, lnk := range m.Paths[focusInfo].ExitILinks {
				exitILinkInputs[lnk] = 1
			}
			result = pather.ProcessRiver(m.GetRiverTop(), exitILinkInputs, focusInfo, supportingInfos, pth, true)
		} else {
			result = pather.ProcessRiver(m.GetRiverTop(), exitILinkInputs, focusInfo, supportingInfos, pth, false)
		}
		print(focusInfo.Uid)
		println(result)
		// Open cascade upon failure of process river
		if !result {
			m.OpenCascade(focusInfo)
		}
		totalCount := 1
		correctCount := 0
		if result {
			correctCount = 1
		}
		// Process all cascades and use correctness to gauge fitgoodness
		cascadesToProcess := m.Cascades[focusInfo]
		for _, cascade := range cascadesToProcess {
			pth.TakeSnapshot()
			cascadeSuccess := pather.ProcessCascadeWithIVariation(cascade, focusInfo, supportingInfos, pth)
			pth.RestoreSnapshot()
			totalCount = totalCount + 1
			if cascadeSuccess {
				correctCount = correctCount + 1
			}
			_ = cascade
		}

		// Update fitgoodness
		var goodness float32
		if len(cascadesToProcess) == 0 {
			goodness = 1.0
		} else {
			print(correctCount)
			print("/")
			println(totalCount)
			goodness = float32(correctCount) / float32(totalCount)
		}
		_ = goodness
		oldFitGoodness := m.InfoFitGoodnesses[focusInfo]
		m.InfoFitGoodnesses[focusInfo] = (7 * oldFitGoodness / 8) + (goodness / 8)
		println(m.InfoFitGoodnesses[focusInfo])
	}
	println()
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
