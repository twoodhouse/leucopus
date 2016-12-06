package memory

import (
	"testing"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/truthTable"
)

func TestMemoryUsage(t *testing.T) {
	mem := New()
	info1 := info.New("a")
	info2 := info.New("b")
	// info3 := info.New("c")

	mem.SetRiver(info1, []int{0, 1, 1, 0, 1, 0, 1, 1, 1})
	mem.SetRiver(info2, []int{1, 0, 1, 1, 0, 1, 0, 1, 1, 1})
	// mem.SetRiver(info3, []int{0, 1, 1, 1})
	mem.PrintRiver()
	mem.OpenCascade(info1, []*info.Info{info2})
	testsConfig := mem.GenerateRiverTests(info1, []*info.Info{info2}, 0)
	tableI1 := truthTable.New([]int{0, 0, 1, 1})
	tableR := truthTable.New([]int{0, 1})
	linkI1 := truthTable.NewLink(tableI1)
	linkR := truthTable.NewLink(tableR)
	sourceLink := truthTable.NewLink(truthTable.NewEntryTable())
	truthTable.AttachLinks(sourceLink, linkI1, 0)
	truthTable.AttachLinks(linkI1, linkI1, 1)
	truthTable.AttachLinks(linkI1, linkR, 0)
	//sourceLink(0) -> linkI1 -> linkR(output) TODO: How do I specify internal memory states?
	// memory.CalculateOutput()

	// sourceLink.Forward()

	_ = testsConfig
	// TestPath(testsConfig[0], info1, []*info.Info{info2}, []*truthTable.Link{})
}
