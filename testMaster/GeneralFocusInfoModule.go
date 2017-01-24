package testMaster

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
)

type GeneralFocusInfoModule struct {
	Mem   *memory.Memory
	Infos []*info.Info
}

func NewGeneralFocusInfoModule(mem *memory.Memory) *GeneralFocusInfoModule {
	infos := make([]*info.Info, 0)
	for info := range mem.River {
		infos = append(infos, info)
	}
	var entity = GeneralFocusInfoModule{
		mem,
		infos,
	}
	return &entity
}

func (fim *GeneralFocusInfoModule) GetFocusInfo() *info.Info {
	worstGoodness := float32(1)
	var associatedInfo *info.Info
	for _, info := range fim.Infos {
		if fim.Mem.InfoFitGoodnesses[info] < worstGoodness {
			worstGoodness = fim.Mem.InfoFitGoodnesses[info]
			associatedInfo = info
		}
	}
	return associatedInfo
}
