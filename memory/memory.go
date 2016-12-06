package memory

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/truthTable"
)

type Memory struct {
	River           map[*info.Info][]int
	Paths           map[*info.Info]*pather.Path
	FitGoodnesses   map[*info.Info]float32
	ExitILinkInputs map[*info.Info]map[*truthTable.Link]int
	Cascades        map[*info.Info][]map[*info.Info][]int
	Depths          map[*info.Info]int
	defaultDepth    int
}

func New() *Memory {
	var entity = Memory{
		make(map[*info.Info][]int),
		make(map[*info.Info]*pather.Path),
		make(map[*info.Info]float32),
		make(map[*info.Info]map[*truthTable.Link]int),
		make(map[*info.Info][]map[*info.Info][]int),
		make(map[*info.Info]int),
		50,
	}
	return &entity
}

func (m *Memory) ProcessNextIteration(values map[*info.Info]int) {
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

//TODO: test this function
func (m *Memory) OpenCascade(nfo *info.Info, supportingInfos []*info.Info) {
	if _, ok := m.Cascades[nfo]; !ok {
		m.Cascades[nfo] = make([]map[*info.Info][]int, 0)
	}
	cascade := make(map[*info.Info][]int)
	for _, supportingInfo := range supportingInfos {
		cascade[supportingInfo] = m.River[supportingInfo] //adds whole related river row to the cascade
	}
	m.Cascades[nfo] = append(m.Cascades[nfo], cascade)
}

//TODO: test this function
func (m *Memory) GenerateAllCascadesTestsGivenInfoAndSupport(nfo *info.Info, supportingInfos []*info.Info, depth int) [][]map[*info.Info][]int {
	allTests := make([][]map[*info.Info][]int, 0)
	nfoCascades := m.Cascades[nfo]
	for _, nfoCascade := range nfoCascades {
		hasAllSupport := true
		for _, supportingInfo := range supportingInfos {
			if _, ok := nfoCascade[supportingInfo]; !ok {
				hasAllSupport = false
			}
		}
		if hasAllSupport {
			result := m.GenerateCascadeTests(nfoCascade, nfo, supportingInfos, depth)
			allTests = append(allTests, result)
		}
	}
	return allTests
}

/*
This function needs to generate tests from a cascade, rather than river
River test (singular) is only 1 test which is run on the latest data using ExitILink input values gained from the last river test.
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
