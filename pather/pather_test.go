//note: this file might not work anymore since memory is the driving class

package pather

import (
	"testing"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/truthTable"
)

func TestPathProcessing1(t *testing.T) {
	println("construction test 1")
	mem := memory.New()
	info1 := info.New("result")
	info2 := info.New("source")

	mem.SetRiver(info1, []int{1, 0, 1, 1, 0, 1, 0, 1, 1})
	mem.SetRiver(info2, []int{1, 0, 1, 1, 0, 1, 0, 1, 1, 1})

	pathr := New([]*info.Info{info1, info2})
	testsConfig := mem.GenerateCascadeTests(mem.River, info1, []*info.Info{info2}, 4)

	pth := NewPath([]*info.Info{info2})
	ILink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)
	for i := 0; i < len(testsConfig); i++ {
		pth.ExitILinks[0].Inputs[0] = 0
		memory.GeneralPrintRiverTool(testsConfig[i])
		result := pathr.ProcessTest(testsConfig[i], info1, []*info.Info{info2}, pth)
		println(result)
	}
	println("")
}

func TestPathProcessing2(t *testing.T) {
	println("construction test 2")
	mem := memory.New()
	info1 := info.New("i1")
	info2 := info.New("i2")
	info3 := info.New("i3")

	mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 1})
	mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pathr := New([]*info.Info{info1, info2, info3})
	testsConfig := mem.GenerateCascadeTests(mem.River, info1, []*info.Info{info2, info3}, 4)

	pth := NewPath([]*info.Info{info2, info3})
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 0, 0, 1}), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)

	for i := 0; i < len(testsConfig); i++ {
		pth.ExitILinks[0].Inputs[0] = 1
		memory.GeneralPrintRiverTool(testsConfig[i])
		result := pathr.ProcessTest(testsConfig[i], info1, []*info.Info{info2, info3}, pth)
		println(result)
	}
	println("")
}

func TestPathExitIVariation(t *testing.T) {
	println("IVariation test 1")
	mem := memory.New()
	info1 := info.New("i1")
	info2 := info.New("i2")
	info3 := info.New("i3")

	mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 1})
	mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pathr := New([]*info.Info{info1, info2, info3})
	testsConfig := mem.GenerateCascadeTests(mem.River, info1, []*info.Info{info2, info3}, 4)

	pth := NewPath([]*info.Info{info2, info3})
	midLink := pth.AddLinkFromLinks(pth.EntryLinks, truthTable.New([]int{0, 0, 0, 1}), false)
	ILink := pth.AddLinkFromLinks([]*truthTable.Link{midLink}, truthTable.New([]int{0, 1}), true)
	truthTable.AttachLinks(ILink, pth.ExitLink, 0)

	for i := 0; i < len(testsConfig); i++ {
		// pth.ExitILinks[0].Inputs[0] = 1
		memory.GeneralPrintRiverTool(testsConfig[i])
		/* TODO: I need to consider more deeply what info I will use after I variation.
		I may want info from just the IVariation function, or perhaps from the regular process one (percent of passes?)
		*/
		result := pathr.ProcessTestWithIVariation(testsConfig[i], info1, []*info.Info{info2, info3}, pth)
		println(result)
	}
	println("")
}

func TestPathProcessingDeep(t *testing.T) {
	println("Deep function construction testing")
	mem := memory.New()
	info1 := info.New("a")
	info2 := info.New("b")

	mem.SetRiver(info1, []int{1, 0, 1, 1, 0, 1, 0, 1, 1})
	mem.SetRiver(info2, []int{1, 0, 1, 1, 0, 1, 0, 1, 1, 1})
	// mem.PrintRiver()
	pathr := New([]*info.Info{info1, info2})

	tableI1 := truthTable.New([]int{0, 1})
	tableR := truthTable.New([]int{0, 1})
	linkI1 := truthTable.NewLink(tableI1, false)
	linkR := truthTable.NewLink(tableR, false)
	sourceLink := truthTable.NewLink(truthTable.NewEntryTable(), false)
	exitLink := truthTable.NewLink(truthTable.NewEntryTable(), false)
	exitILink := truthTable.NewLink(truthTable.NewEntryTable(), true)
	exitILink.Inputs[0] = 1
	truthTable.AttachLinks(sourceLink, linkI1, 0)
	truthTable.AttachLinks(linkI1, exitILink, 0)
	truthTable.AttachLinks(exitILink, linkR, 0)
	truthTable.AttachLinks(linkR, exitLink, 0)

	testsConfig := mem.GenerateCascadeTests(mem.River, info1, []*info.Info{info2}, 1)

	// for _, e := range testsConfig {
	// 	memory.GeneralPrintRiverTool(e)
	// }

	sourceLinkAssociation := make(map[*info.Info]*truthTable.Link)
	sourceLinkAssociation[info2] = sourceLink //this associates info2 to the sourceLink/Table so the process can correctly assign the input

	// defaultConfig := []int{0} //implement this soon
	// memory.GeneralPrintRiverTool(testsConfig[0])

	for i := 0; i < len(testsConfig); i++ {
		result := pathr.ProcessTest_Deep(testsConfig[i], info1, []*info.Info{info2}, sourceLinkAssociation, exitLink, []*truthTable.Link{exitILink})
		// _ = result
		println(result)

		// println("source")
		// sourceLink.Print()
		// println("I1")
		// linkI1.Print()
		// println("R")
		// linkR.Print()
		// println("exit")
		// exitLink.Print()
		// println("exitI")
		// exitILink.Print()
	} //Next TODO: simplify construction of paths with some pather utilities. This will also allow more testing, which is necessary.
}
