package utility

import (
	"sort"

	"github.com/twoodhouse/leucopus/info"
)

func GetNextDynamicCombination(infoNums []int, totalNum int) []int {
	if len(infoNums) == 0 {
		return []int{0}
	}

	if len(infoNums) == totalNum {
		return infoNums
	}

	//TODO: finish this (maybe rewrite)
	nextInfoNums := []int{}

	//while ...
	movedOne := false

	carryCount := 0
	farthestBack := totalNum
	hitEmpty := false
	for i := totalNum - 1; !hitEmpty && i >= 0; i-- {
		acted := false
		for _, infoNum := range infoNums {
			if i == infoNum {
				carryCount = carryCount + 1
				acted = true
				farthestBack = i
			}
		}
		if !acted {
			hitEmpty = true
		}
	}

	if carryCount > 0 {
		setPoint := 0
		//find point to set out number of infoNums
		found := false
		for i := farthestBack - 1; i >= 0 && !found; i-- {
			for _, infoNum := range infoNums {
				if infoNum == i {
					found = true
					setPoint = i + 1
				}
			}
		}
		// print(setPoint)
		for i := setPoint; i < carryCount+setPoint+1; i++ {
			nextInfoNums = append(nextInfoNums, i)
		}
		for _, num := range infoNums {
			if num < setPoint-1 {
				nextInfoNums = append(nextInfoNums, num)
			}
		}
	} else {
		for i := totalNum - 2; i >= 0 && !movedOne; i-- { //no use checking the last element, thus the -2
			for _, infoNum := range infoNums {
				//does the infoNums contain an element at the next slot?
				if infoNum == i {
					containsNext := false
					for _, infoNum2 := range infoNums {
						// print(infoNum2)
						// println("!")
						if infoNum2 == i+1 {
							containsNext = true
						}
					}
					if !containsNext {
						//simple move is available
						nextInfoNums = append(nextInfoNums, i+1)
						for _, num := range infoNums {
							if num < i {
								nextInfoNums = append(nextInfoNums, num)
							}
						}
						movedOne = true
					}
				}
			}
		}
	}

	return nextInfoNums
}

func GetFullDynamicStringCombinations(lst []string) [][]string { //modify this
	var fullCombinations [][]string
	lastIntComb := []int{}
	for len(lastIntComb) < len(lst) {
		lastIntComb = GetNextDynamicCombination(lastIntComb, len(lst))
		currStringComb := []string{}
		for _, el := range lastIntComb {
			currStringComb = append(currStringComb, lst[el])
		}
		fullCombinations = append(fullCombinations, currStringComb)
	}
	return fullCombinations
}

func GetUidFromInfos(infos []*info.Info) string {
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
