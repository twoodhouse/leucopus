package testMaster

import (
	"math"
	"sort"

	"github.com/cznic/mathutil"
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/truthTable"
)

type GeneralPathModule struct {
	Mem                *memory.Memory
	LastPathCollection map[*info.Info]map[string]*pather.Path
}

func NewGeneralPathModule(mem *memory.Memory) *GeneralPathModule {
	var entity = GeneralPathModule{
		mem,
		make(map[*info.Info]map[string]*pather.Path),
	}
	return &entity
}

func (fim *GeneralPathModule) GetPath(nfo *info.Info, supportingInfos []*info.Info) *pather.Path {
	//TODO: current comments below
	//current issue: How can I effectively iterate over the full functionality
	//perhaps there is an easy iteration method (my original one - or similar), which can do it without too much loss

	//****** - increment the lowest modRow (and waterfall the process upward if it is the last available)
	// I1: A, B, A-B
	// I2: A, B, I1, A-B, A-I1, B-I1, A-B-I1 (EXCEPT the options containing something the option before chose and not an additional I node)
	// R: I2<unlinked I and input nodes>, I2<unlinked I and input nodes>-(full permutation of elements not necessary in R)
	// In order to do this algorithm, the available I2 and R options will need to be determined
	//
	/*
		Example 1:
		I1: <A, B, A-B> = A
		R: <1-B, 1-B*(A)> = 1-B

		Example 2:
		I1: <A, B, A-B> = A
		I2: <B, 1, A-1, B-1, A-B-1> = B
		R: <2-1, 2-1(A, B, A-B)>

		Example 3:
		I1: <A, B, A-B> = A
		I2: <B, 1, A-1, B-1, A-B-1> = A-1
		R: <2-B, 2-B(A, 1, A-1)> = 2-B-A <- is A allowed here? ***YES***

		Trevor's rules of pathing
		1. A complete mental model for an output may be composed of any number of I nodes and a single R node
		2. I nodes are ordered
		3. Any number of inputs may be assigned to an I node, or an R node, from available sources
		4. For I nodes, available sources include the looped output of previously assigned I nodes, or any info point.
		5. Info points are designated by their name ("a", "B", etc) while I nodes are designated by their number.
		6. Available sources are controlled for I nodes according to the following rules
			a. An I node may not contain any portion of input to a previous I node ("Z") unless it also contains an I input which is not referenced by "Z"
			b. The R node must contain as inputs any I nodes or Info points which are not referenced by an I node, and optionally any I nodes or Info points
	*/

	//first get last path
	var lastPath *pather.Path
	if _, ok := fim.LastPathCollection[nfo]; ok {
		if _, ok := fim.LastPathCollection[nfo][getUidFromInfos(supportingInfos)]; ok {
			lastPath = fim.LastPathCollection[nfo][getUidFromInfos(supportingInfos)]
		} else {
			lastPath = nil
		}
	} else {
		fim.LastPathCollection[nfo] = make(map[string]*pather.Path)
		lastPath = nil
	}

	//now get dynamic string combination
	sInfoNames := []string{}
	for _, supportingInfo := range supportingInfos {
		sInfoNames = append(sInfoNames, supportingInfo.Uid)
	}
	dynamicStringCombinations := utility.GetFullDynamicStringCombinations(sInfoNames)
	_ = dynamicStringCombinations

	//now get number of I nodes
	supportingUid := getUidFromInfos(supportingInfos)
	var numINodes int
	if lastPath != nil {
		numINodes = len(lastPath.ExitILinks)
	} else {
		numINodes = 0
	}
	_ = numINodes

	//now get max node size
	var maxNodeSize int
	if lastPath != nil {
		maxNodeSize = getLargestInputNum(lastPath)
	} else {
		maxNodeSize = 0
	}
	_ = maxNodeSize

	//edge case for first path
	if lastPath == nil {
		fim.LastPathCollection[nfo] = make(map[string]*pather.Path)
		t := []int{}
		for i := 0; i < 2; i++ {
			t = append(t, 2)
		}
		pth := pather.NewPath(getInfosFromUids([]string{dynamicStringCombinations[0][0]}, supportingInfos))
		midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New(t), false)
		truthTable.AttachLinks(midLink, pth.ExitLink, 0)
		fim.LastPathCollection[nfo][supportingUid] = pth
		return pth
	}

	//make combination map and find current location in the combination map
	cMap := make([][][]string, 0)
	cMapRow0 := make([][]string, len(dynamicStringCombinations))
	copy(cMapRow0, dynamicStringCombinations)
	cMap = append(cMap, cMapRow0)
	//Loop: (TODO)
	//	determine location in previous cMapRow
	//	make combination map for this next row
	//end loop

	//remember to have it check for lower max node Sizes in other parts of the map (how do I do this?)

	//**NOTE: the following section is un-important. Just for testing.
	//edge case 1: no existing path used last time
	if _, ok := fim.LastPathCollection[nfo]; !ok {
		fim.LastPathCollection[nfo] = make(map[string]*pather.Path)
		t := []int{}
		for i := 0; i < int(math.Pow(2, float64(len(supportingInfos)))); i++ {
			t = append(t, 2)
		}
		pth := pather.NewPath(supportingInfos)
		midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New(t), false)
		truthTable.AttachLinks(midLink, pth.ExitLink, 0)
		fim.LastPathCollection[nfo][supportingUid] = pth
		return pth
	}

	t := []int{}
	for i := 0; i < int(math.Pow(2, float64(len(supportingInfos)))); i++ {
		t = append(t, 2)
	}
	pth := pather.NewPath(supportingInfos)
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New(t), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)
	return pth
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
