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
	MiddleLinks     []*truthTable.Link
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
		println()
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
		println()
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
		println()
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
	println()
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
		return false
	}
	return true
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
	for i := 0; i < int(math.Exp2(float64(numExitILinks))); i++ {
		binaryStrRow := strings.Split(strconv.FormatInt(int64(i), 2), "")
		binaryIntRow := make([]int, numExitILinks)
		for i, e := range binaryStrRow {
			binaryIntRow[i], _ = strconv.Atoi(e)
		}
		for j, e := range binaryIntRow {
			pth.ExitILinks[j].Inputs[0] = e
			result := ProcessTest(test, nfo, supportingInfos, pth)
			if result {
				return true
			}
		}
	}
	return false
}

func ProcessTest(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path) bool {
	return ProcessTest_Deep(test, nfo, supportingInfos, pth.LinkAssociation, pth.ExitLink, pth.ExitILinks)
}

func ProcessTest_Deep(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, sourceLinkAssociation map[*info.Info]*truthTable.Link, exitLink *truthTable.Link, exitILinks []*truthTable.Link) bool {

	for i := 0; i < len(test[nfo]); i++ {
		for _, exitILink := range exitILinks {
			exitILink.Process()
			exitILink.Forward()
		}

		var unfinishedLinks []*truthTable.Link
		for _, supportingInfo := range supportingInfos {
			if _, ok := sourceLinkAssociation[supportingInfo]; ok {
				history := test[supportingInfo]
				sourceLinkAssociation[supportingInfo].Output = history[i]
				intermediateUnfinishedLinks := sourceLinkAssociation[supportingInfo].Forward() //forward should return a list of links at which the process stopped
				for _, newLink := range intermediateUnfinishedLinks {                          //now add all the unifinished links together
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
		if exitLink.Output != test[nfo][i] {
			return false
		}
	}
	return true
}
