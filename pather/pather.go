package pather

import (
	"math"
	"strconv"
	"strings"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/truthTable"
)

var uidCounter = 0

type Path struct {
	LinkAssociation map[*info.Info]*truthTable.Link
	MiddleLinks     []*truthTable.Link //The last middle link is R!!!
	EntryLinks      []*truthTable.Link
	ExitILinks      []*truthTable.Link
	ExitLink        *truthTable.Link
	Age             int
	Uid             string
	EntryInfos      []*info.Info
}

// func (p *Path) Copy() *Path {
// 	var newPath *Path
// 	newPath = NewPath(p.EntryInfos)
// 	return newPath
// }

func (p *Path) Print() {
	print("Entry Info Uids ")
	for _, nfo := range p.EntryInfos {
		print(nfo.Uid)
		print(", ")
	}
	println()
	println("************ Entry Links  ************")
	for _, lnk := range p.EntryLinks {
		print("-")
		print("\t\t")
		print(lnk.Uid)
		print("\t\t")
		for _, underTargetLink := range lnk.TargetLinks {
			print(underTargetLink.Uid)
			print(", ")
		}
		print("\t")
		lnk.Table.Print()
	}
	println("************ Middle Links ************")
	for _, lnk := range p.MiddleLinks {
		for _, underSourceLink := range lnk.SourceLinks {
			print(underSourceLink.Uid)
			print(", ")
		}
		print("\t\t")
		print(lnk.Uid)
		print("\t\t")
		for _, underTargetLink := range lnk.TargetLinks {
			print(underTargetLink.Uid)
			print(", ")
		}
		print("\t")
		lnk.Table.Print()
	}
	println("************ Exit I Links ************")
	for _, lnk := range p.ExitILinks {
		for _, underSourceLink := range lnk.SourceLinks {
			print(underSourceLink.Uid)
			print(", ")
		}
		print("\t\t")
		print(lnk.Uid)
		print("\t\t")
		for _, underTargetLink := range lnk.TargetLinks {
			print(underTargetLink.Uid)
			print(", ")
		}
		print("\t")
		lnk.Table.Print()
	}
	println("************  Exit Link   ************")
	lnk := p.ExitLink
	for _, underSourceLink := range lnk.SourceLinks {
		print(underSourceLink.Uid)
		print(", ")
	}
	print("\t\t")
	print(lnk.Uid)
	print("\t\t")
	for _, underTargetLink := range lnk.TargetLinks {
		print(underTargetLink.Uid)
		print(", ")
	}
	print("\t")
	lnk.Table.Print()
}

func NewPath(entryInfos []*info.Info) *Path {
	exitLink := truthTable.NewLink(truthTable.NewEntryTable(), false)
	entryNum := len(entryInfos)
	var entryLinks []*truthTable.Link
	for i := 0; i < entryNum; i++ {
		entryLinks = append(entryLinks, truthTable.NewLink(truthTable.NewEntryTable(), false))
	}
	var exitILinks []*truthTable.Link
	linkAssociation := make(map[*info.Info]*truthTable.Link, 0)
	for i, entryInfo := range entryInfos {
		linkAssociation[entryInfo] = entryLinks[i]
	}
	uid := strconv.Itoa(uidCounter)
	uidCounter = uidCounter + 1
	var entity = Path{
		linkAssociation,
		[]*truthTable.Link{},
		entryLinks,
		exitILinks,
		exitLink,
		0,
		uid,
		entryInfos,
	}
	return &entity
}

func (p *Path) TakeSnapshot() {
	for _, middleLink := range p.MiddleLinks {
		middleLink.TakeSnapshot()
	}
	for _, entryLink := range p.EntryLinks {
		entryLink.TakeSnapshot()
	}
	for _, exitILink := range p.ExitILinks {
		exitILink.TakeSnapshot()
	}
	p.ExitLink.TakeSnapshot()
}

func (p *Path) TakeStaticTablesSnapshot() {
	for _, middleLink := range p.MiddleLinks {
		if middleLink != p.MiddleLinks[len(p.MiddleLinks)-1] {
			middleLink.Table.TakeSnapshot()
		}
	}
}

func (p *Path) TakeRTableSnapshot() {
	p.MiddleLinks[len(p.MiddleLinks)-1].Table.TakeSnapshot()
}

func (p *Path) RestoreRTableSnapshot() {
	p.MiddleLinks[len(p.MiddleLinks)-1].Table.RestoreSnapshot()
}

func (p *Path) TakeRTableSnapshot_lower() {
	p.MiddleLinks[len(p.MiddleLinks)-1].Table.TakeSnapshot_lower()
}

func (p *Path) RestoreRTableSnapshot_lower() {
	p.MiddleLinks[len(p.MiddleLinks)-1].Table.RestoreSnapshot_lower()
}

func (p *Path) RestoreStaticTablesSnapshot() {
	for _, middleLink := range p.MiddleLinks {
		if middleLink != p.MiddleLinks[len(p.MiddleLinks)-1] {
			middleLink.Table.RestoreSnapshot()
		}
	}
}

func (p *Path) RestoreSnapshot() {
	for _, middleLink := range p.MiddleLinks {
		middleLink.RestoreSnapshot()
	}
	for _, entryLink := range p.EntryLinks {
		entryLink.RestoreSnapshot()
	}
	for _, exitILink := range p.ExitILinks {
		exitILink.RestoreSnapshot()
	}
	p.ExitLink.RestoreSnapshot()
}

func (p *Path) AgePath() {
	p.Age = p.Age + 1
}

func (p *Path) AddLinkFromLinks(sourceLinks []*truthTable.Link, table *truthTable.TruthTable, isExitILink bool) *truthTable.Link {
	newLink := truthTable.NewLink(table, isExitILink)
	for i, sourceLink := range sourceLinks {
		truthTable.AttachLinks(sourceLink, newLink, i)
	}
	if isExitILink {
		p.ExitILinks = append(p.ExitILinks, newLink)
	} else {
		p.MiddleLinks = append(p.MiddleLinks, newLink)
	}
	return newLink
}

func ProcessRiver(mostRecent map[*info.Info]int, exitILinkInputs map[*truthTable.Link]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path, setInitialExitILinks bool) bool {
	if setInitialExitILinks {
		for _, exitILink := range pth.ExitILinks {
			exitILink.Inputs[0] = exitILinkInputs[exitILink]
		}
	}
	for _, exitILink := range pth.ExitILinks {
		exitILink.Process()
		exitILink.Forward()
	}
	var unfinishedLinks []*truthTable.Link
	for _, supportingInfo := range supportingInfos {
		if _, ok := pth.LinkAssociation[supportingInfo]; ok {
			history := mostRecent[supportingInfo]
			pth.LinkAssociation[supportingInfo].Output = history
			intermediateUnfinishedLinks := pth.LinkAssociation[supportingInfo].Forward() //forward should return a list of links at which the process stopped
			for _, newLink := range intermediateUnfinishedLinks {                        //now add all the unifinished links together
				existsAlready := false
				for _, refLink := range unfinishedLinks {
					if newLink == refLink {
						existsAlready = true
					}
				}
				if !existsAlready {
					unfinishedLinks = append(unfinishedLinks, newLink)
				}
			}
		}
	}
	if pth.ExitLink.Output != mostRecent[nfo] {
		// println(pth.ExitLink.Output)
		// for k, v := range mostRecent {
		// 	print(k.Uid)
		// 	print(":")
		// 	println(v)
		// }
		// pth.MiddleLinks[0].Print()
		return false
	}
	return true
}

func ProcessRiver2(mostRecent map[*info.Info]int, exitILinkInputs map[*truthTable.Link]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path, setInitialExitILinks bool) (bool, int) {
	if setInitialExitILinks {
		for _, exitILink := range pth.ExitILinks {
			exitILink.Inputs[0] = exitILinkInputs[exitILink]
		}
	}
	//configure ExitILinks so that the initial targets of them have the appropriate values
	for _, exitILink := range pth.ExitILinks {
		exitILink.Process()
		exitILink.Forward()
	}

	//for each supportingInfo of the focused Path, Forward to recursively dive through ALL the rest of the path truth tables (R likely last)
	for _, supportingInfo := range supportingInfos {
		if _, ok := pth.LinkAssociation[supportingInfo]; ok {
			history := mostRecent[supportingInfo]
			pth.LinkAssociation[supportingInfo].Output = history
			pth.LinkAssociation[supportingInfo].Forward()
		}
	}
	pth.AssumeBestOfR(mostRecent[nfo])

	// println(nfo.Uid)
	// println(pth.ExitLink.Output)
	// println(mostRecent[nfo])
	// println()

	var index int
	index = pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.LastProcessLocation

	if pth.ExitLink.Output != mostRecent[nfo] {
		return false, index
	}
	return true, index
}

func (pth *Path) AssumeBestOfR(desiredResult int) {
	if pth.MiddleLinks[len(pth.MiddleLinks)-1].Output == 2 {
		//update the truth table output to be the expected result
		pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.ReplaceInputValue(pth.MiddleLinks[len(pth.MiddleLinks)-1].Inputs, desiredResult)
		//update ExitLink (and middle link?) output to be the expected result
		pth.MiddleLinks[len(pth.MiddleLinks)-1].Output = desiredResult
		pth.ExitLink.Output = desiredResult
	}
}

func ProcessAllCascadesTogether(cascades []map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path, balance []int) int { //return the maximum number that were correct
	pth.TakeSnapshot()             //remember, this snapshot will not store or reset the truthtable values.
	pth.TakeStaticTablesSnapshot() // this snapshot will store the STATIC truthtable values
	pth.TakeRTableSnapshot()
	/*	Quite simply, I want all the cascades which pass so that they must be consistent with one another and the main river.
		This means that I can let ProcessCascadeWithIVariation deal with the iteration over entry values.
		Iteration over static truthtable values should happen at the top level: here, because it must be done the same accross all cascades
		Assuming best of R should be done at the low level and the river balance recorded across a whole static truthtable combo at this level.
		Can assuming best of R be done on values which R already determined? Yes. ???
		Everything ought to be reset when changing static I tables. Record best case continually so that it can be returned
		 and so that the the static tables and R table can be set from the stored values at the end.
	*/
	var totalNumStaticOutputs int
	for _, middleLink := range pth.MiddleLinks {
		if middleLink != pth.MiddleLinks[len(pth.MiddleLinks)-1] {
			totalNumStaticOutputs = totalNumStaticOutputs + len(middleLink.Table.Outputs) //TODO: update this to be only 2s? shouldn
		}
	}
	bestNumCorrect := -1
	bestBinaryIntRow := []int{}
	bestBalanceLocal := []int{}
	for i := 0; i < len(balance); i++ {
		bestBalanceLocal = append(bestBalanceLocal, 0)
	}
	println()

	for i := 0; i < int(math.Exp2(float64(totalNumStaticOutputs))); i++ { //next iterate over possible ExitINode initial values
		balanceLocal := []int{}
		for j := 0; j < len(balance); j++ {
			balanceLocal = append(balanceLocal, 0)
		}
		binaryStr := strconv.FormatInt(int64(i), 2)
		intermediate := ""
		for j := 0; j < totalNumStaticOutputs-len(binaryStr); j++ {
			intermediate = intermediate + "0"
		}
		binaryStr = intermediate + binaryStr
		binaryStrRow := strings.Split(binaryStr, "")
		binaryIntRow := make([]int, len(binaryStrRow))
		for i, e := range binaryStrRow {
			binaryIntRow[i], _ = strconv.Atoi(e)
		}
		if nfo.Uid == "i1" { //for testing only
			for _, e := range binaryStrRow {
				print(e)
			}
			println()
		}
		counter := 0
		for _, middleLink := range pth.MiddleLinks {
			if middleLink != pth.MiddleLinks[len(pth.MiddleLinks)-1] {
				//now assign outputs to the static tables appropriately using the counter
				for i, _ := range middleLink.Table.Outputs {
					middleLink.Table.Outputs[i] = binaryIntRow[counter]
					counter = counter + 1
				}
			}
		}
		//static tables should be assigned appropriately at this point. Now process cascade with entry variation.
		numCorrect := 0
		for _, cascade := range cascades {
			//reset R table outputs to 2
			for i, _ := range pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs {
				pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs[i] = 2
			}
			success := ProcessCascadeWithIVariation(cascade, nfo, supportingInfos, pth)
			if success {
				if nfo.Uid == "i1" {
					println("success")
				}
				//if success, possibly the R table has been modified. Must keep track of which parameters changed in the river balance
				for _, outputVal_old := range pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.OutputsSnapshot {
					for i, outputVal := range pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs {
						if outputVal != outputVal_old {
							// this output value changed from the baseline over the course of this cascade
							if outputVal == 1 {
								balanceLocal[i] = balanceLocal[i] + 1
							} else if outputVal == 0 {
								balanceLocal[i] = balanceLocal[i] - 1
							}
						}
					}
				}
				numCorrect = numCorrect + 1
			} else {
				if nfo.Uid == "i1" {
					for k, v := range cascade { //for testing only
						print(k.Uid)
						print(": ")
						for _, e := range v {
							print(e)
						}
						println()
					}
				}
			}
			pth.RestoreRTableSnapshot() // no need to keep modified R table between cascades. It will be recorded in the river balance
		}
		if numCorrect > bestNumCorrect {
			if nfo.Uid == "i1" {
				println("found best!")
				println(numCorrect)
				println(len(cascades))
			}
			bestNumCorrect = numCorrect
			bestBinaryIntRow = binaryIntRow
			bestBalanceLocal = balanceLocal
		}
	}
	pth.RestoreStaticTablesSnapshot() // this snapshot will store the STATIC truthtable values

	// now REASSIGN the best set of static table outputs
	counter := 0
	for _, middleLink := range pth.MiddleLinks {
		if middleLink != pth.MiddleLinks[len(pth.MiddleLinks)-1] {
			//now assign outputs to the static tables appropriately using the counter
			for i, _ := range middleLink.Table.Outputs {
				middleLink.Table.Outputs[i] = bestBinaryIntRow[counter]
				counter = counter + 1
			}
		}
	}
	//now make R table output modifications based on the best balance (local)
	if len(cascades) > 0 {
		for i, balanceOutput := range bestBalanceLocal {
			if balanceOutput > 0 {
				println("setting high")
				pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs[i] = 1
			} else if balanceOutput < 0 {
				println("setting low")
				pth.MiddleLinks[len(pth.MiddleLinks)-1].Table.Outputs[i] = 0
			}
		}
	}
	_ = bestBalanceLocal
	_ = bestBinaryIntRow
	pth.RestoreSnapshot()
	return bestNumCorrect
}

func ProcessCascadeWithIVariation(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path) bool {
	//print cascade for testing
	// print("test for ")
	// println(nfo.Uid)
	// for k, v := range test {
	// 	print(k.Uid)
	// 	print(":")
	// 	for _, e := range v {
	// 		print(e)
	// 		print(",")
	// 	}
	// 	println()
	// }
	numExitILinks := len(pth.ExitILinks)
	if numExitILinks == 0 {
		result := ProcessTest(test, nfo, supportingInfos, pth)
		if result {
			return true
		} else {
			return false
		}
	} else {
		for i := 0; i < int(math.Exp2(float64(numExitILinks))); i++ { //next iterate over possible ExitINode initial values
			binaryStr := strconv.FormatInt(int64(i), 2)
			intermediate := ""
			for i := 0; i < numExitILinks-len(binaryStr); i++ {
				intermediate = intermediate + "0"
			}
			binaryStr = intermediate + binaryStr
			binaryStrRow := strings.Split(binaryStr, "")
			binaryIntRow := make([]int, len(binaryStrRow))
			for i, e := range binaryStrRow {
				binaryIntRow[i], _ = strconv.Atoi(e)
			}
			for j, e := range binaryIntRow {
				pth.ExitILinks[j].Inputs[0] = e
			}
			pth.TakeRTableSnapshot_lower() //this snapshot will store the R truthtable values
			if nfo.Uid == "i1" {
				pth.Print()
			}
			result := ProcessTest(test, nfo, supportingInfos, pth)
			if result {
				return true //remember cascades are only pass or fail
			}
			pth.RestoreRTableSnapshot_lower()
		}
		return false
	}
}

func ProcessTest(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path) bool {
	return ProcessTest_Deep(test, nfo, supportingInfos, pth.LinkAssociation, pth.ExitLink, pth.ExitILinks, pth)
}

func ProcessTest_Deep(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, sourceLinkAssociation map[*info.Info]*truthTable.Link, exitLink *truthTable.Link, exitILinks []*truthTable.Link, pth *Path) bool {

	for i := 0; i < len(test[nfo]); i++ {
		for _, exitILink := range exitILinks {
			exitILink.Process()
			exitILink.Forward()
		}

		// var unfinishedLinks []*truthTable.Link
		for _, supportingInfo := range supportingInfos {
			if _, ok := sourceLinkAssociation[supportingInfo]; ok {
				history := test[supportingInfo]
				sourceLinkAssociation[supportingInfo].Output = history[i]
				sourceLinkAssociation[supportingInfo].Forward()
				// intermediateUnfinishedLinks := sourceLinkAssociation[supportingInfo].Forward() //forward should return a list of links at which the process stopped
				// for _, newLink := range intermediateUnfinishedLinks {                          //now add all the unifinished links together
				// 	existsAlready := false
				// 	for _, refLink := range unfinishedLinks {
				// 		if newLink == refLink {
				// 			existsAlready = true
				// 		}
				// 	}
				// 	if !existsAlready {
				// 		unfinishedLinks = append(unfinishedLinks, newLink)
				// 	}
				// }
			}
		}
		pth.AssumeBestOfR(test[nfo][i])

		if exitLink.Output != test[nfo][i] {
			return false
		}
	}
	return true
}

func ProcessTest_Deep_old(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, sourceLinkAssociation map[*info.Info]*truthTable.Link, exitLink *truthTable.Link, exitILinks []*truthTable.Link) bool {

	for i := 0; i < len(test[nfo]); i++ {
		for _, exitILink := range exitILinks {
			exitILink.Process()
			exitILink.Forward()
		}

		// var unfinishedLinks []*truthTable.Link
		for _, supportingInfo := range supportingInfos {
			if _, ok := sourceLinkAssociation[supportingInfo]; ok {
				history := test[supportingInfo]
				sourceLinkAssociation[supportingInfo].Output = history[i]
				sourceLinkAssociation[supportingInfo].Forward()
				// intermediateUnfinishedLinks := sourceLinkAssociation[supportingInfo].Forward() //forward should return a list of links at which the process stopped
				// for _, newLink := range intermediateUnfinishedLinks {                          //now add all the unifinished links together
				// 	existsAlready := false
				// 	for _, refLink := range unfinishedLinks {
				// 		if newLink == refLink {
				// 			existsAlready = true
				// 		}
				// 	}
				// 	if !existsAlready {
				// 		unfinishedLinks = append(unfinishedLinks, newLink)
				// 	}
				// }
			}
		}
		//TODO: add R truthtable checkpoint?
		//TODO: add AssumeBest here?

		if exitLink.Output != test[nfo][i] {
			return false
		}
	}
	return true
}
