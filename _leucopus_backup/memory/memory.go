package memory

import "github.com/twoodhouse/leucopus/info"

type Memory struct {
	River        map[*info.Info][]int
	Cascades     map[*info.Info]map[*info.Info][]int
	Depths       map[*info.Info]int
	defaultDepth int
}

func New() *Memory {
	var entity = Memory{
		make(map[*info.Info][]int),
		make(map[*info.Info]map[*info.Info][]int),
		make(map[*info.Info]int),
		50,
	}
	return &entity
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

func (m *Memory) SetRiver(nfo *info.Info, vals []int) {
	// if memory does not contain this info, add a depth for it
	if _, ok := m.Depths[nfo]; !ok {
		m.Depths[nfo] = m.defaultDepth
	}
	m.River[nfo] = vals
}

func (m *Memory) OpenCascade(nfo *info.Info, supportingInfos []*info.Info) { //TODO: ???make this able to determine supportingInfos on its own
	for _, supportingInfo := range supportingInfos {
		m.Cascades[nfo] = make(map[*info.Info][]int)
		m.Cascades[nfo][supportingInfo] = m.River[supportingInfo] //adds whole related river row to the cascade
	}
}

func (m *Memory) GenerateRiverTests(nfo *info.Info, supportingInfos []*info.Info, depth int) []map[*info.Info][]int {
	tests := make([]map[*info.Info][]int, 0)
	riverDry := false
	shortestLength := len(m.River[nfo])
	for _, supportingInfo := range supportingInfos {
		if len(m.River[supportingInfo]) < shortestLength {
			shortestLength = len(m.River[supportingInfo])
		}
	}
	for i := shortestLength; !riverDry; i-- { //TODO: fix this section to take from end of list to the beginning
		riverTest := make(map[*info.Info][]int, 0)
		offset := len(m.River[nfo]) - shortestLength
		riverTest[nfo] = m.River[nfo][i-depth-1+offset : i+offset]
		for _, supportingInfo := range supportingInfos {
			offset := len(m.River[supportingInfo]) - shortestLength
			riverTest[supportingInfo] = m.River[supportingInfo][i-depth-1+offset : i+offset]
			if i-depth-1 <= 0 {
				riverDry = true
			}
		}
		// m.PrintRiverTool(riverTest)
		tests = append(tests, riverTest)
		riverTest = nil
	}
	// for i := 0; !riverDry; i++ { //TODO: fix this section to take from end of list to the beginning
	// 	riverTest := make(map[*info.Info][]int, 0)
	// 	riverTest[nfo] = m.River[nfo][i : i+depth+1]
	// 	for _, supportingInfo := range supportingInfos {
	// 		riverTest[supportingInfo] = m.River[supportingInfo][i : i+depth+1]
	// 		if len(m.River[supportingInfo]) == i+depth+1 || len(m.River[nfo]) == i+depth+1 {
	// 			riverDry = true
	// 		}
	// 	}
	// }
	// intermediate node values are determined in the Tester.
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
