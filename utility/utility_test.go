package utility

import "testing"

func TestUtility(t *testing.T) {
	// e := []int{}
	// for i := 0; i < 15; i++ {
	// 	e = GetNextDynamicCombination(e, 4)
	// 	for _, el := range e {
	// 		print(el)
	// 	}
	// 	println()
	// }
}

func TestMetaPath(t *testing.T) {
	metaPath := NewMetaPath([]string{"A", "B"}, 2)
	metaPath.Explore()
}
