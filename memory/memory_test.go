package memory

import (
	"testing"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/truthTable"
)

func TestPathPrint(t *testing.T) {
	println("Path print test")
	mem := New()
	info1 := info.New("i1")
	info2 := info.New("i2")
	info3 := info.New("i3")

	mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 0})
	mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pth := pather.NewPath([]*info.Info{info2, info3})
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 0, 0, 1}), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)

	pth.Print()
	println()
}

func TestPathRiverProcess(t *testing.T) {
	println("River process test 1")
	mem := New()
	info1 := info.New("i1")
	info2 := info.New("i2")
	info3 := info.New("i3")

	mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 0})
	mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pathr := pather.New([]*info.Info{info1, info2, info3})

	pth := pather.NewPath([]*info.Info{info2, info3})
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 0, 0, 1}), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)

	//setup memory as the parent dictator class will do
	mem.Paths[info1] = pth
	mem.ExitILinkInputs[info1] = make(map[*truthTable.Link]int)
	mem.ExitILinkInputs[info1][ILink] = 1

	//now run actual test
	result := pathr.ProcessRiver(mem.GetRiverTop(), mem.ExitILinkInputs[info1], info1, []*info.Info{info2, info3}, pth)

	for _, exitILink := range pth.ExitILinks {
		mem.ExitILinkInputs[info1][exitILink] = exitILink.Inputs[0]
	}
	println(result)

	newRow := make(map[*info.Info]int)
	newRow[info1] = 0
	newRow[info2] = 1
	newRow[info3] = 1
	mem.ProcessNextIteration(newRow)

	result = pathr.ProcessRiver(mem.GetRiverTop(), mem.ExitILinkInputs[info1], info1, []*info.Info{info2, info3}, pth)
	println(result)

	println("")
}

func TestPathExitIVariation(t *testing.T) {
	println("Cascade IVariation test 1")
	mem := New()
	info1 := info.New("i1")
	info2 := info.New("i2")
	info3 := info.New("i3")

	mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 1})
	mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pathr := pather.New([]*info.Info{info1, info2, info3})
	testsConfig := mem.GenerateCascadeTests(mem.River, info1, []*info.Info{info2, info3}, 4)

	pth := pather.NewPath([]*info.Info{info2, info3})
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 0, 0, 1}), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)

	for i := 0; i < len(testsConfig); i++ {
		// pth.ExitILinks[0].Inputs[0] = 1
		GeneralPrintRiverTool(testsConfig[i])
		/* TODO: I need to consider more deeply what info I will use after I variation.
		I may want info from just the IVariation function, or perhaps from the regular process one (percent of passes?)
		*/
		result := pathr.ProcessCascadeWithIVariation(testsConfig[i], info1, []*info.Info{info2, info3}, pth)
		println(result)
	}
	println("")
}
