package pather

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/truthTable"
)

type Path struct {
	LinkAssociation map[*info.Info]*truthTable.Link
	MiddleLinks     []*truthTable.Link
	EntryLinks      []*truthTable.Link
	ExitILinks      []*truthTable.Link
	ExitLink        *truthTable.Link
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
	var entity = Path{
		linkAssociation,
		[]*truthTable.Link{},
		entryLinks,
		exitILinks,
		exitLink,
	}
	return &entity
}

func (p *Path) AddLinkFromLinks(sourceLinks []*truthTable.Link, table *truthTable.TruthTable, isExitILink bool) *truthTable.Link {
	newLink := truthTable.NewLink(table, isExitILink)
	for i, sourceLink := range sourceLinks {
		truthTable.AttachLinks(sourceLink, newLink, i)
	}
	if isExitILink {
		p.ExitILinks = append(p.ExitILinks, newLink)
	}
	return newLink
}

type Pather struct {
	memory *memory.Memory
	infos  []*info.Info
}

func New(memory *memory.Memory, infos []*info.Info) *Pather {
	var entity = Pather{
		memory,
		infos,
	}
	return &entity
}

func (p *Pather) ProcessTest(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, pth *Path) bool {
	return p.ProcessTest_Deep(test, nfo, supportingInfos, pth.LinkAssociation, pth.ExitLink, pth.ExitILinks)
}

func (p *Pather) ProcessTest_Deep(test map[*info.Info][]int, nfo *info.Info, supportingInfos []*info.Info, sourceLinkAssociation map[*info.Info]*truthTable.Link, exitLink *truthTable.Link, exitILinks []*truthTable.Link) bool {

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
