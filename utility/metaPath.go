package utility

import "strconv"

type MetaPath struct {
	AllSources         []string
	TopMetaRow         *MetaRow
	NumINodes          int
	IncreaseMaxRowSize bool
}

func NewMetaPath(sources []string, numINodes int) *MetaPath {
	var metaRow *MetaRow
	originalSources := make([]string, len(sources))
	copy(originalSources, sources)
	// if numINodes > 0 {
	// 	sources = append([]string{"I1"}, sources...)
	// }
	for i := numINodes; i > 0; i-- { //should this traverse the other direction?
		st := "I" + strconv.Itoa(i)
		sources = append(sources, st)
	}
	inputSet := GetFullDynamicStringCombinations(sources)
	inputSetWithoutI1 := make([][]string, 0)
	if numINodes == 0 {
		inputSetWithoutI1 = append(inputSetWithoutI1, sources)
	} else {
		for _, input := range inputSet {
			if len(input) != 1 || input[0] != "I1" {
				inputSetWithoutI1 = append(inputSetWithoutI1, input)
			}
		}
	}

	metaRow = NewMetaRow(nil, inputSetWithoutI1, nil, nil, originalSources, numINodes)
	//now load the first metaPath to test:
	mrTest := metaRow
	for mrTest.FinalChoice == nil {
		mrTest = mrTest.Delve(mrTest.Choices[0], false)
	}
	var entity = MetaPath{
		sources,
		metaRow,
		numINodes,
		false,
	}
	return &entity
}

func (mp *MetaPath) GetMaxRowSize() int {
	max := 1
	mrTest := mp.TopMetaRow
	for mrTest.Child != nil {
		mrTest = mrTest.Child
		if len(mrTest.LeadChoice) > max {
			max = len(mrTest.LeadChoice)
		}
	}
	if len(mrTest.FinalChoice) > max {
		max = len(mrTest.FinalChoice)
	}
	return max
}

func (mp *MetaPath) GetDeepestRow() *MetaRow {
	mrTest := mp.TopMetaRow
	for mrTest.Child != nil {
		mrTest = mrTest.Child
	}
	return mrTest
}

// func (mp *MetaPath) Improve2() {
// 	//row size functionality. To increase, set the variable
// 	maxRowSize := mp.GetMaxRowSize()
// 	if mp.IncreaseMaxRowSize {
// 		maxRowSize = maxRowSize + 1
// 		mp.IncreaseMaxRowSize = false
// 	}
// 	//set deepest row
// 	deepestRow := mp.GetDeepestRow()
// 	//loop until a change is made
// 	changeMade := false
// 	curRow := deepestRow //start the row pointer on the deepest row
// 	for !changeMade {
//
// 	}
// }

func (mp *MetaPath) Improve2() bool {
	if isFullyImproved(mp.GetDeepestRow()) {
		return false
	}
	maxRowSize := mp.GetMaxRowSize()
	if mp.IncreaseMaxRowSize {
		maxRowSize = maxRowSize + 1
		mp.IncreaseMaxRowSize = false
	}
	deepestRow := mp.GetDeepestRow()
	var activeRow *MetaRow
	activeRow = deepestRow
	madeChange := false
	for !madeChange {
		if activeRow.Child == nil {
			//on final row
			activeRowChoice := activeRow.FinalChoice
			var locationInRow int
			for i, choice := range activeRow.Choices {
				if compareInputLists(choice, activeRowChoice) {
					locationInRow = i
				}
			}
			if locationInRow < len(activeRow.Choices)-1 && len(activeRow.Choices[locationInRow+1]) <= maxRowSize {
				activeRow.Delve(activeRow.Choices[locationInRow+1], false)
				madeChange = true
			} else {
				activeRow = activeRow.Parent
			}
		} else {
			//not on final row
			activeRowChoice := activeRow.Child.LeadChoice
			var locationInRow int
			for i, choice := range activeRow.Choices {
				if compareInputLists(choice, activeRowChoice) {
					locationInRow = i
				}
			}
			if locationInRow < len(activeRow.Choices)-1 && len(activeRow.Choices[locationInRow+1]) <= maxRowSize {
				activeRow := activeRow.Delve(activeRow.Choices[locationInRow+1], false)
				firstRowSuccess := DelveFirstOption_Conditional(activeRow, maxRowSize, mp)
				if firstRowSuccess {
					madeChange = true
				} else {
				}
			} else {
				if activeRow.Parent != nil {
					activeRow = activeRow.Parent
				} else {
					mp.IncreaseMaxRowSize = true
					maxRowSize = maxRowSize + 1
					activeRow = deepestRow
				}
			}
		}
	}
	return true
}

func isFullyImproved(deepestRow *MetaRow) bool {
	mrTest := deepestRow
	if !compareInputLists_noOrder(mrTest.FinalChoice, mrTest.Choices[len(mrTest.Choices)-1]) {
		return false
	}
	for mrTest.Parent != nil {
		if !compareInputLists_noOrder(mrTest.LeadChoice, mrTest.Parent.Choices[len(mrTest.Parent.Choices)-1]) {
			return false
		}
		mrTest = mrTest.Parent
	}
	return true
}

//
// func (mp *MetaPath) Improve() {
// 	maxRowSize := mp.GetMaxRowSize()
// 	if mp.IncreaseMaxRowSize {
// 		maxRowSize = maxRowSize + 1
// 		mp.IncreaseMaxRowSize = false
// 	}
// 	deepestRow := mp.GetDeepestRow()
// 	var curFollowingRow *MetaRow
// 	changeMade := false
// 	var curChoiceOverride []string
// 	for !changeMade {
// 		var curChoice []string
// 		var curRow *MetaRow
// 		if curFollowingRow == nil {
// 			if curChoiceOverride != nil {
// 				curChoice = curChoiceOverride
// 				curChoiceOverride = nil
// 			} else {
// 				curChoice = deepestRow.FinalChoice
// 			}
// 			curRow = deepestRow
// 		} else {
// 			if curChoiceOverride != nil {
// 				curChoice = curChoiceOverride
// 			} else {
// 				curChoice = curFollowingRow.LeadChoice
// 			}
// 			curRow = curFollowingRow.Parent
// 		}
// 		// find current location so we can move on to the next
// 		// print(curRow)
// 		// print(" ")
// 		// for _, e := range curChoice {
// 		// 	print(e)
// 		// }
// 		// println()
//
// 		var locationInRow int
// 		for i, choice := range curRow.Choices {
// 			if compareInputLists(choice, curChoice) {
// 				locationInRow = i
// 			}
// 		}
// 		//case 1: move on to next option
// 		if (len(curRow.Choices) > locationInRow+1 && len(curRow.Choices[locationInRow+1]) <= maxRowSize) || curChoiceOverride != nil {
// 			returnedRow := curRow.Delve(curRow.Choices[locationInRow+1], false)
// 			if curRow.Name != "R" {
// 				firstRowSuccess := DelveFirstOption_Conditional(returnedRow, maxRowSize)
// 				if firstRowSuccess {
// 					curChoiceOverride = nil
// 					changeMade = true
// 				} else {
// 				}
// 			} else {
// 				changeMade = true
// 			}
// 			//remember to break or put extra conditional on next if statement
// 		}
//
// 		if curChoiceOverride != nil {
// 			returnedRow := curRow.Delve(curRow.Choices[locationInRow], false)
// 			firstRowSuccess := DelveFirstOption_Conditional(returnedRow, maxRowSize)
// 			if firstRowSuccess {
// 				curChoiceOverride = nil
// 				changeMade = true
// 			} else {
// 			}
// 		}
//
// 		//case 2: case 1 does not find an appropriate option. Move up a row and try again.
// 		// if curRow.Parent == nil && len(curRow.Child.LeadChoice) == maxRowSize && IsFinalRowSizePoint(curRow, maxRowSize) { //OLD: Does this work?
// 		// curRow.PrintChoices()
//
// 		if curRow.Parent == nil && IsFinalRowSizePoint(curRow, maxRowSize) { //increase row size
// 			mp.IncreaseMaxRowSize = true
// 			curFollowingRow = nil
// 			curChoiceOverride = nil
// 			deepestRow = mp.GetDeepestRow()
// 			curFollowingRow = nil
// 			curChoiceOverride = nil
// 			println("increase size")
// 		} else if IsFinalRowChoiceForMax(curRow, curChoice, maxRowSize) && curRow.Parent != nil { //go up a row
// 			curFollowingRow = curRow
// 			curChoiceOverride = nil
// 			println("moving up")
// 		} else {
// 			if _, ok := curRow.Choices[locationInRow+1]; ok {
// 				curChoiceOverride = curRow.Choices[locationInRow+1]
// 				println("override")
// 			} else {
// 				mp.IncreaseMaxRowSize = true
// 			}
// 		}
// 	}
// }

func IsFinalRowChoiceForMax(curRow *MetaRow, choice []string, max int) bool {
	found := false
	for _, ch := range curRow.Choices {
		if found == true {
			if len(ch) > max {
				return true
			}
			return false
		}
		if compareInputLists(choice, ch) {
			found = true
		}
	}
	println("the choice was not found in the current row")
	return false
}

func IsFinalRowSizePoint(mrTop *MetaRow, max int) bool {
	mrTest := mrTop
	foundSmaller := false
	for mrTest.Child != nil {
		if !compareInputLists(mrTest.Choices[len(mrTest.Choices)-1], mrTest.Child.LeadChoice) {
			foundSmaller = true
		}
		mrTest = mrTest.Child
	}
	if !compareInputLists(mrTest.Choices[len(mrTest.Choices)-1], mrTest.FinalChoice) {
		foundSmaller = true
	}
	return !foundSmaller
}

func DelveFirstOption_Conditional(mr *MetaRow, max int, mp *MetaPath) bool {
	mrTest := mr
	for mrTest.Name != "R" {
		mrTest = mrTest.Delve(mrTest.Choices[0], false)
	}
	mrTest.Delve(mrTest.Choices[0], false)
	bottomRow := mrTest
	topRow := mrTest
	for topRow.Parent != nil {
		topRow = topRow.Parent
	}
	_ = bottomRow
	//now I have the bottom and the top. Iterate up till a fair combo is found
	currentRow := bottomRow
	currentRowNum := 0
	indexes := []int{}
	for mp.GetMaxRowSize() > max {
		// print(mr.Name)
		// print("-")
		// mp.Print()
		if len(indexes) == currentRowNum {
			indexes = append(indexes, 0)
		}
		if currentRow == bottomRow {
			//at bottom
			if len(currentRow.Choices[indexes[currentRowNum]+1]) <= max {
				currentRow.Delve(currentRow.Choices[indexes[currentRowNum]+1], false)
			} else {
				if currentRow == mr.Parent {
					return false
				}
				//failure, look at next element
				if indexes[currentRowNum] < len(currentRow.Choices)-2 {
					indexes[currentRowNum] = indexes[currentRowNum] + 1
				} else {
					// if at end of element list, go up a level
					currentRowNum = currentRowNum + 1
					currentRow = currentRow.Parent
					if currentRow == mr.Parent {
						return false
					}
					if len(indexes) == currentRowNum {
						indexes = append(indexes, 0)
					}
					indexes[currentRowNum] = indexes[currentRowNum] + 1
				}
			}
		} else {
			//not at bottom
			cont := true
			if len(currentRow.Choices[indexes[currentRowNum]+1]) <= max {
				currentAltRow := currentRow.Delve(currentRow.Choices[indexes[currentRowNum]+1], false)
				// println("going deeper")
				// print(mr.Name)
				// print("-")
				// mp.Print()
				firstRowSuccess := DelveFirstOption_Conditional(currentAltRow, max, mp)
				if firstRowSuccess {
					return true
				} else {
				}
			}
			if cont == true {
				//failure, look at next element
				if currentRow == mr.Parent {
					return false
				}
				if indexes[currentRowNum] < len(currentRow.Choices)-2 {
					indexes[currentRowNum] = indexes[currentRowNum] + 1
				} else {
					// if at end of element list, go up a level
					currentRowNum = currentRowNum + 1
					currentRow = currentRow.Parent
					if currentRow == mr.Parent {
						return false
					}
					if len(indexes) == currentRowNum {
						indexes = append(indexes, 0)
					}
					indexes[currentRowNum] = indexes[currentRowNum] + 1
				}
			}
		}
	}
	return true
}

func (mp *MetaPath) GetChosen() [][]string {
	mrTest := mp.TopMetaRow
	for mrTest.Child != nil {
		mrTest = mrTest.Child
	}
	canGoDeeper := true
	choices := make([][]string, 0)
	if mrTest.FinalChoice == nil {
		println("Final choice is not set. This GetChosen() may have been run too early.")
	} else {
		choices = append(choices, mrTest.FinalChoice)
	}
	for canGoDeeper {
		if mrTest.Parent == nil {
			canGoDeeper = false
		} else {
			choices = append([][]string{mrTest.LeadChoice}, choices...)
			mrTest = mrTest.Parent
		}
	}
	return choices
}

func (mp *MetaPath) Print() {
	mrTest := mp.TopMetaRow
	for mrTest.Child != nil {
		mrTest = mrTest.Child
	}
	mrTest.PrintChoices()
}

type MetaRow struct {
	Name          string
	LeadChoice    []string
	Choices       [][]string
	OriginChoices [][]string
	Parent        *MetaRow
	Child         *MetaRow
	Sources       []string
	numINodes     int
	FinalChoice   []string
}

func NewMetaRow(leadChoice []string, choices [][]string, originChoices [][]string, parent *MetaRow, sources []string, numINodes int) *MetaRow {
	if originChoices == nil {
		originChoices = choices
	}
	name := ""

	if parent == nil {
		if numINodes == 0 {
			name = "R"
		} else {
			name = "I1"
		}
	} else {
		val, _ := strconv.Atoi(parent.Name[len(parent.Name)-1:])
		if val == numINodes {
			name = "R"
		} else {
			name = "I" + strconv.Itoa(val+1)
		}
	}

	var entity = MetaRow{
		name,
		leadChoice,
		choices,
		originChoices,
		parent,
		nil,
		sources,
		numINodes,
		nil,
	}
	return &entity
}

func (mr *MetaRow) Delve(choice []string, test bool) *MetaRow {
	/*
		This is where all the real important code is. This is where the laws of pathing come into play
	*/
	if mr.Name == "R" {
		mr.FinalChoice = choice
		return mr
	} else {

		//First Level: General Availability
		depth, _ := strconv.Atoi(mr.Name[len(mr.Name)-1:])
		inputs := make([]string, 0)
		mod := 0
		if depth != mr.numINodes {
			mod = 1
		}
		for i := mr.numINodes - 1 + mod; i > 0; i-- { //should this traverse the other direction?
			st := "I" + strconv.Itoa(i)
			inputs = append(inputs, st)
		}
		for _, choice := range mr.Sources {
			inputs = append(inputs, choice)
		}
		gac := GetFullDynamicStringCombinations(inputs) //generally available combinations

		//Now enforce limitations on the list of inputs
		limitedGac := make([][]string, 0)
		if depth != mr.numINodes { //I case
			//Limits: not just self or any previous entry without additional I
			for _, permutation := range gac {
				pDepth, _ := strconv.Atoi(permutation[0][len(permutation[0])-1:])
				if !compareInputLists(permutation, choice) && !comparePreviousInputListsPlusMod(permutation, mr, choice) && !(len(permutation) == 1 && pDepth == depth+1) { // not existing adjacent choice
					limitedGac = append(limitedGac, permutation)
				}
			}
		} else { //R case
			//get leadChoices (and current choice) used in all metarows
			usedNodes := make([]string, 0)
			for _, e := range choice {
				if e != mr.Name {
					usedNodes = addIfNotPresent(e, usedNodes)
				}
			}
			canGoDeeper := true
			mrTest := mr
			for canGoDeeper {
				for _, e := range mrTest.LeadChoice {
					if e[:1] == "I" {
						nameDepth, _ := strconv.Atoi(mrTest.Name[1:])
						nameDepth = nameDepth - 1
						eDepth, _ := strconv.Atoi(e[1:])
						// println(nameDepth)
						// println(eDepth)
						if nameDepth > eDepth {
							usedNodes = addIfNotPresent(e, usedNodes)
						}
					} else {
						usedNodes = addIfNotPresent(e, usedNodes)
					}
				}
				if mrTest.Parent != nil {
					mrTest = mrTest.Parent
				} else {
					canGoDeeper = false
				}
			}
			//get full permutations of all I nodes and inputs
			rac := make([][]string, 0) //r available combinations
			rInputs := make([]string, 0)

			for _, source := range mr.Sources {
				rInputs = append(rInputs, source)
			}
			for i := 1; i < mr.numINodes+1; i++ {
				val := strconv.Itoa(i)
				rInputs = append(rInputs, "I"+val)
			}

			notUsedNodes := make([]string, 0)
			for _, input := range rInputs {
				found := false
				for _, usedInput := range usedNodes {
					if input == usedInput {
						found = true
					}
				}
				if !found {
					notUsedNodes = append(notUsedNodes, input)
				}
			}

			rac = GetFullDynamicStringCombinations(rInputs)
			//distribute used set into full permutation set
			limitedGac = append(limitedGac, notUsedNodes)
			for _, pSet := range rac {
				for _, input := range notUsedNodes {
					pSet = addIfNotPresent(input, pSet)
				}
				limitedGac = addIfNotPresent_Deeper(pSet, limitedGac)
			}
		}

		newMr := NewMetaRow(choice, limitedGac, mr.OriginChoices, mr, mr.Sources, mr.numINodes)
		if !test {
			mr.Child = newMr
		}
		return newMr
	}
}

func (mr *MetaRow) PrintChoices() {
	canGoDeeper := true
	mrTest := mr
	choices := make([][]string, 0)
	if mr.FinalChoice == nil {
		println("Final choice is not set. This print may have been run too early.")
	} else {
		choices = append(choices, mr.FinalChoice)
	}
	for canGoDeeper {
		if mrTest.Parent == nil {
			canGoDeeper = false
		} else {
			choices = append([][]string{mrTest.LeadChoice}, choices...)
			mrTest = mrTest.Parent
		}
	}
	print("{")
	for _, choice := range choices {
		for _, el := range choice {
			print(el)
		}
		print(", ")
	}
	println("}")
}

func addIfNotPresent(el string, lst []string) []string {
	present := false
	for _, e := range lst {
		if e == el {
			present = true
		}
	}
	if !present {
		lst = append(lst, el)
	}
	return lst
}

func addIfNotPresent_Deeper(el []string, lst [][]string) [][]string {
	present := false
	for _, e := range lst {
		if compareInputLists_noOrder(e, el) {
			present = true
		}
	}
	if !present {
		lst = append(lst, el)
	}
	return lst
}

func elContainedIn(el string, lst []string) bool {
	for _, lstEl := range lst {
		if el == lstEl {
			return false
		}
	}
	return true
}

func comparePreviousInputListsPlusMod(la []string, mr *MetaRow, choice []string) bool {
	if allAInBCompareInputLists_noAdditionalI(choice, la) { // la is the permutation
		return true
	}
	// depth, _ := strconv.Atoi(mr.Name[len(mr.Name)-1:]) //TODO: fix this
	// println(depth)
	// if len(la) == 1 && la[0] == "I"+strconv.Itoa(depth) {
	// 	return true
	// }
	canGoDeeper := true
	mrTest := mr
	for canGoDeeper {
		if compareInputLists(mrTest.LeadChoice, la) {
			return true
		}
		if len(mrTest.LeadChoice) > 0 && allAInBCompareInputLists_noAdditionalI(mrTest.LeadChoice, la) {
			return true
		}
		if mrTest.Parent != nil {
			mrTest = mrTest.Parent
		} else {
			canGoDeeper = false
		}
	}
	return false
}

func compareInputLists(la []string, lb []string) bool {
	if len(la) != len(lb) {
		return false
	}
	for i, _ := range la {
		if la[i] != lb[i] {
			return false
		}
	}
	return true
}

func compareInputLists_noOrder(la []string, lb []string) bool {
	if len(la) != len(lb) {
		return false
	}
	for _, ea := range la {
		found := false
		for _, eb := range lb {
			if eb == ea {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func allAInBCompareInputLists_noAdditionalI(la []string, lb []string) bool {
	for _, aEl := range la {
		inB := false
		for _, bEl := range lb {
			if aEl == bEl {
				inB = true
			}
		}
		if !inB {
			return false
		}
	}
	//so then all of A is in B
	//now if B also has an I node which is not in A, return false
	for _, bEl := range lb {
		if bEl[0:1] == "I" {
			inA := false
			for _, aEl := range la {
				if aEl == bEl {
					inA = true
				}
			}
			if !inA {
				return false
			}
		}
	}
	return true
}

//
// func (mp *MetaPath) Explore() {
// 	currentMetaRow := mp.TopMetaRow
// 	var parents []*MetaRow
// 	parentTestRow := currentMetaRow
// 	for parentTestRow.Parent != nil {
// 		parents = append(parents, parentTestRow.Parent)
// 		parentTestRow = parentTestRow.Parent
// 	}
// 	for i := len(parents) - 1; i >= 0; i++ {
// 		for j := 0; j < i; j++ {
// 			print("\t")
// 		}
// 		println(parents[i].Name)
// 	}
// 	for i := 0; i < len(parents); i++ {
// 		print("\t")
// 	}
// 	for _, el := range currentMetaRow.Choices {
// 		for _, e := range el {
// 			print(e)
// 		}
// 		print(", ")
// 	}
// 	println()
// 	print("-> ")
// }
