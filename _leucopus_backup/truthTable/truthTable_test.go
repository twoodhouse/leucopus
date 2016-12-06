package truthTable

import (
	"testing"
)

func TestNodeCreation(t *testing.T) {
	truthTable := New([]int{2, 2, 1, 0, 0, 1, 0, 0})
	// truthTable := New([]int{1, 2, 0, 2})
	// println(truthTable.Process([]int{1, 1}))
	truthTable.Fill()
	truthTable.Print()
	// truthTable.Print()
}
