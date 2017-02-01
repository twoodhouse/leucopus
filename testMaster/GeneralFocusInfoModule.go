package testMaster

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
)

type GeneralFocusInfoModule struct {
	Mem       *memory.Memory
	Infos     []*info.Info
	LastIndex int
}

func NewGeneralFocusInfoModule(mem *memory.Memory) *GeneralFocusInfoModule {
	var entity = GeneralFocusInfoModule{
		mem,
		[]*info.Info{},
		0,
	}
	return &entity
}

func (fim *GeneralFocusInfoModule) PullRiverInfos() {
	for nfo, _ := range fim.Mem.River {
		fim.Infos = append(fim.Infos, nfo)
	}
}

func (fim *GeneralFocusInfoModule) GetFocusInfo() *info.Info {
	newIndex := fim.LastIndex + 1
	if newIndex == len(fim.Infos) {
		newIndex = 0
	}
	choice := fim.Infos[newIndex]
	fim.LastIndex = newIndex
	return choice
}

// func (fim *GeneralFocusInfoModule) GetFocusInfo() *info.Info {
// 	worstGoodness := float32(1)
// 	var associatedInfo *info.Info
// 	for _, info := range fim.Infos {
// 		if fim.Mem.InfoFitGoodnesses[info] < worstGoodness {
// 			worstGoodness = fim.Mem.InfoFitGoodnesses[info]
// 			associatedInfo = info
// 		}
// 	}
// 	return associatedInfo
// }
