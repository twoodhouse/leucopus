package testMaster

import (
	"testing"

	"github.com/twoodhouse/leucopus/info"
)

func TestStringCombination(t *testing.T) {
	result := getFullDynamicStringCombinations([]string{"a", "b", "c", "d"})
	for _, set := range result {
		for _, e := range set {
			print(e)
		}
		println()
	}
}

func TestTestMasterUse(t *testing.T) {
	println("testing testMaster creation and use 1")
	tm := New()

	info1 := info.New("a")
	info2 := info.New("b")
	info3 := info.New("c")

	tm.Mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	tm.Mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 0})
	tm.Mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	for i := 0; i < 5; i++ {
		// tm.GoopNextPath() //use the modules to get next path for testing. This may also do partial fill for the new path from existing paths
	}

	// links := pth.GetAllLinks()
	// pth, pathGoodness := tm.GoopPath(pth, links) //iterate over various truth table possibilities for the links (only tests replacing 2s)

}

func TestPather(t *testing.T) {
	println("testing Path Module with case 1")
	tm := New()

	info1 := info.New("a")
	info2 := info.New("b")
	info3 := info.New("c")

	tm.Mem.SetRiver(info1, []int{1, 0, 0, 0, 1, 0, 0, 1})
	tm.Mem.SetRiver(info2, []int{1, 1, 0, 0, 1, 1, 0, 1, 0})
	tm.Mem.SetRiver(info3, []int{1, 0, 1, 1, 1, 0, 1, 1, 1})

	pathModule := NewGeneralPathModule(tm.Mem)
	pth := pathModule.GetPath(info1, []*info.Info{info2, info3})
	pth.Print()
	pth = pathModule.GetPath(info1, []*info.Info{info2, info3})
	pth.Print()
}

// func TestNextSupportingInfo(t *testing.T) {
// 	println("testing permutation thing")
// 	result := getNextInfoNums([]int{2, 3, 4}, 5)
// 	println("result")
// 	for _, e := range result {
// 		println(e)
// 	}
// }
