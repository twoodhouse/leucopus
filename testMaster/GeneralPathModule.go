package testMaster

import (
	"math"
	"sort"
	"strconv"

	"github.com/cznic/mathutil"
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/truthTable"
	"github.com/twoodhouse/leucopus/utility"
)

type GeneralPathModule struct {
	Mem                    *memory.Memory
	LastMetaPathCollection map[*info.Info]map[string]*utility.MetaPath
}

func NewGeneralPathModule(mem *memory.Memory) *GeneralPathModule {
	var entity = GeneralPathModule{
		mem,
		make(map[*info.Info]map[string]*utility.MetaPath),
	}
	return &entity
}

func (fim *GeneralPathModule) GetPath(nfo *info.Info, supportingInfos []*info.Info) *pather.Path {
	//get last path
	var lastMetaPath *utility.MetaPath
	if _, ok := fim.LastMetaPathCollection[nfo]; ok {
		if _, ok := fim.LastMetaPathCollection[nfo][getUidFromInfos(supportingInfos)]; ok {
			lastMetaPath = fim.LastMetaPathCollection[nfo][getUidFromInfos(supportingInfos)]
		} else {
			lastMetaPath = nil
		}
	} else {
		fim.LastMetaPathCollection[nfo] = make(map[string]*utility.MetaPath)
		lastMetaPath = nil
	}

	//determine list of supportingInfo names
	supportingInfoNames := []string{}
	_ = supportingInfoNames
	for _, nfo := range supportingInfos {
		supportingInfoNames = append(supportingInfoNames, nfo.Uid)
	}

	var metaPath *utility.MetaPath
	if lastMetaPath == nil {
		metaPath = utility.NewMetaPath(supportingInfoNames, 0)
	} else {
		metaPath = lastMetaPath
		improveSuccess := metaPath.Improve2()
		if !improveSuccess {
			//determine current depth
			metaPath = utility.NewMetaPath(supportingInfoNames, lastMetaPath.NumINodes+1)
		}
	}
	fim.LastMetaPathCollection[nfo][getUidFromInfos(supportingInfos)] = metaPath

	chosen := metaPath.GetChosen()
	// for _, el := range chosen {
	// 	for _, e := range el {
	// 		print(e)
	// 	}
	// 	print(",")
	// }
	// println()

	// now make path, from array of inputs for each iNode
	pth := pather.NewPath(supportingInfos)
	//make all INodes first so they can be linked from
	ILinks := []*truthTable.Link{}
	for i := 0; i < len(chosen)-1; i++ {
		ILink := pth.AddLinkFromLinks(nil, truthTable.New([]int{0, 1}), true)
		ILinks = append(ILinks, ILink)
	}
	//make middle nodes
	for choiceNum, choice := range chosen {
		//locate correct links to link from infosToLink and I values
		linksToLink := make([]*truthTable.Link, 0)
		for _, infoName := range choice {
			for i, supportingInfo := range supportingInfos {
				if infoName == supportingInfo.Uid {
					linksToLink = append(linksToLink, pth.EntryLinks[i])
				}
			}
			if infoName[0:1] == "I" {
				val, _ := strconv.Atoi(infoName[1:])
				linksToLink = append(linksToLink, ILinks[val-1])
			}
		}
		//make links
		midLink := pth.AddLinkFromLinks(linksToLink, truthTable.New(getTwosArrayForTableSize(len(linksToLink))), false)
		_ = midLink
		//make necessary modifications to associated ILink
		if choiceNum == len(chosen)-1 {
			truthTable.AttachLinks(midLink, pth.ExitLink, 0)
		} else {
			truthTable.AttachLinks(midLink, ILinks[choiceNum], 0)
		}
	}
	return pth
	//**NOTE: the following section is un-important. Just for testing.
	//edge case 1: no existing path used last time
	// if _, ok := fim.LastPathCollection[nfo]; !ok {
	// 	fim.LastPathCollection[nfo] = make(map[string]*pather.Path)
	// 	t := []int{}
	// 	for i := 0; i < int(math.Pow(2, float64(len(supportingInfos)))); i++ {
	// 		t = append(t, 2)
	// 	}
	// 	pth := pather.NewPath(supportingInfos)
	// 	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New(t), false)
	// 	truthTable.AttachLinks(midLink, pth.ExitLink, 0)
	// 	fim.LastPathCollection[nfo][supportingUid] = pth
	// 	return pth
	// }
	//
	// t := []int{}
	// for i := 0; i < int(math.Pow(2, float64(len(supportingInfos)))); i++ {
	// 	t = append(t, 2)
	// }
	// pth := pather.NewPath(supportingInfos)
	// midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New(t), false)
	// ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	// truthTable.AttachLinks(ILink, pth.ExitLink, 0)
	// return pth
}

func getTwosArrayForTableSize(tableInputNum int) []int {
	t := []int{}
	for i := 0; i < int(math.Pow(2, float64(tableInputNum))); i++ {
		t = append(t, 2)
	}
	return t
}

func getDifficultyFactor(pth *pather.Path) int {
	total := 1
	for _, lnk := range pth.MiddleLinks {
		total = total * int(math.Pow(2, math.Pow(2, float64(len(lnk.SourceLinks)))))
	}
	return total
}

func getLargestInputNum(pth *pather.Path) int {
	max := 0
	for _, lnk := range pth.MiddleLinks {
		if len(lnk.SourceLinks) > max {
			max = len(lnk.SourceLinks)
		}
	}
	return max
}

type ByName []string

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i] < a[j] }

func getInfosFromUids(uids []string, infos []*info.Info) []*info.Info {
	rInfos := make([]*info.Info, 0)
	for _, uid := range uids {
		for _, nfo := range infos {
			if nfo.Uid == uid {
				rInfos = append(rInfos, nfo)
			}
		}
	}
	return rInfos
}

func getSetNumPermutations(lst []string) [][]string { //remove this
	var perms [][]string
	mathutil.PermutationFirst(ByName(lst))
	nLst := make([]string, len(lst))
	copy(nLst, lst)
	perms = append(perms, nLst)
	goOn := true
	for goOn {
		goOn = mathutil.PermutationNext(ByName(lst))
		if goOn {
			nLst := make([]string, len(lst))
			copy(nLst, lst)
			perms = append(perms, nLst)
		}
	}
	return perms
}

func getUidFromInfos(infos []*info.Info) string {
	outList := []string{}
	for _, e := range infos {
		outList = append(outList, e.Uid)
	}
	sort.Strings(outList)
	var finSt string
	for _, e := range outList {
		finSt = finSt + e
	}
	return finSt
}
