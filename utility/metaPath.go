package utility

import (
	"bufio"
	"fmt"
	"os"
)

type MetaPath struct {
	AllSources []string
	TopMetaRow *MetaRow
	numINodes  int
}

func NewMetaPath(sources []string, numINodes int) *MetaPath {
	var metaRow *MetaRow
	if numINodes > 1 {
		metaRow = NewMetaRow("I1", GetFullDynamicStringCombinations(sources), nil)
	} else {
		metaRow = NewMetaRow("R", GetFullDynamicStringCombinations(sources), nil)
	}
	var entity = MetaPath{
		sources,
		metaRow,
		numINodes,
	}
	return &entity
}

type MetaRow struct {
	Name    string
	Choices [][]string
	Parent  *MetaRow
}

func NewMetaRow(source string, choices [][]string, parent *MetaRow) *MetaRow {
	var entity = MetaRow{
		source,
		choices,
		parent,
	}
	return &entity
}

func (mr *MetaRow) Delve(choice string) *MetaRow {

	return NewMetaRow(choice, nil, mr.Parent)
}

func (mp *MetaPath) Explore() {
	currentMetaRow := mp.TopMetaRow
	for {
		var parents []*MetaRow
		parentTestRow := currentMetaRow
		for parentTestRow.Parent != nil {
			parents = append(parents, parentTestRow.Parent)
			parentTestRow = parentTestRow.Parent
		}
		for i := len(parents) - 1; i >= 0; i++ {
			for j := 0; j < i; j++ {
				print("\t")
			}
			println(parents[i].Name)
		}
		for i := 0; i < len(parents); i++ {
			print("\t")
		}
		for _, el := range currentMetaRow.Choices {
			for _, e := range el {
				print(e)
			}
			print(", ")
		}
		println()
		print("-> ")
		scanner := bufio.NewScanner(os.Stdin)
		var text string
		for text != "q" { // break the loop if text == "q"
			fmt.Print("Enter your text: ")
			scanner.Scan()
			text = scanner.Text()
			if text != "q" {
				fmt.Println("Your text was: ", text)
			}
		}

		// if response == ".." && currentMetaRow.Parent != nil {
		// 	currentMetaRow = currentMetaRow.Parent
		// }
		println("ending loop")
	}
}
