package testMaster

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
)

type GeneralSupportingInfosModule struct {
	Mem                   *memory.Memory
	Infos                 []*info.Info
	LastSupportingInfoSet map[*info.Info][]*info.Info
	TriesOnSet            map[*info.Info]int
	MaxTries              map[*info.Info]int
}

func NewGeneralSupportingInfosModule(mem *memory.Memory) *GeneralSupportingInfosModule {
	infos := make([]*info.Info, 0)
	for info := range mem.River {
		infos = append(infos, info)
	}
	var entity = GeneralSupportingInfosModule{
		mem,
		[]*info.Info{},
		make(map[*info.Info][]*info.Info),
		make(map[*info.Info]int),
		make(map[*info.Info]int),
	}
	return &entity
}

func (fim *GeneralSupportingInfosModule) PullRiverInfos() {
	MAX_TRIES := 2
	for nfo, _ := range fim.Mem.River {
		fim.Infos = append(fim.Infos, nfo)
		fim.MaxTries[nfo] = MAX_TRIES
		fim.TriesOnSet[nfo] = MAX_TRIES
	}
}

func (fim *GeneralSupportingInfosModule) GetSupportingInfos(nfo *info.Info) []*info.Info {
	if fim.TriesOnSet[nfo] < fim.MaxTries[nfo]-1 {
		fim.TriesOnSet[nfo] = fim.TriesOnSet[nfo] + 1
		return fim.LastSupportingInfoSet[nfo]
	}
	fim.TriesOnSet[nfo] = 0
	//populate infos from river
	if len(fim.Infos) == 0 {
		for nfo := range fim.Mem.River {
			fim.Infos = append(fim.Infos, nfo)
		}
	}

	if _, ok := fim.LastSupportingInfoSet[nfo]; !ok {
		fim.LastSupportingInfoSet[nfo] = []*info.Info{}
	}
	lastSet := fim.LastSupportingInfoSet[nfo]
	infoNums := []int{}
	for _, supportingInfo := range lastSet {
		for i, stInfo := range fim.Infos {
			if stInfo == supportingInfo {
				infoNums = append(infoNums, i)
			}
		}
	}
	//now I have infoNums correct, so I can calculate the next infoNums
	nextInfoNums := getNextInfoNums(infoNums, len(fim.Infos))

	supportingInfos := []*info.Info{}
	for _, num := range nextInfoNums {
		for i, stInfo := range fim.Infos {
			if num == i {
				supportingInfos = append(supportingInfos, stInfo)
			}
		}
	}
	fim.LastSupportingInfoSet[nfo] = supportingInfos
	return supportingInfos
}

func getNextInfoNums(infoNums []int, totalNum int) []int {
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

func (fim *GeneralSupportingInfosModule) GetExtraInfos(nfo *info.Info) []*info.Info {
	return []*info.Info{}
}
