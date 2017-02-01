package testMaster

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/pather"
)

type TestMaster struct {
	Mem                       *memory.Memory
	FocusInfoModuleUsed       FocusInfoModule
	SupportingInfosModuleUsed SupportingInfosModule
	PathModuleUsed            PathModule
	LastFocus                 *info.Info
	LastSupportingInfos       []*info.Info
	LastPath                  *pather.Path
}

type FocusInfoModule interface {
	GetFocusInfo() *info.Info
	PullRiverInfos()
}

type SupportingInfosModule interface {
	GetSupportingInfos(*info.Info) []*info.Info
	PullRiverInfos()
	GetExtraInfos(*info.Info) []*info.Info
}

type PathModule interface {
	GetPath(*info.Info, []*info.Info) *pather.Path
}

func New() *TestMaster {
	mem := memory.New()
	var entity = TestMaster{
		mem,
		NewGeneralFocusInfoModule(mem),
		NewGeneralSupportingInfosModule(mem),
		NewGeneralPathModule(mem),
		nil,
		nil,
		nil,
	}
	return &entity
}

func (tm *TestMaster) GetNextPath() *pather.Path {
	focusInfo := tm.FocusInfoModuleUsed.GetFocusInfo()
	supportingInfos := tm.SupportingInfosModuleUsed.GetSupportingInfos(focusInfo)
	extraInfos := tm.SupportingInfosModuleUsed.GetExtraInfos(focusInfo) //TODO: this has something to do with cascades?
	//TODO: extraInfos need to be added to the cascade if the goop fails
	pth := tm.PathModuleUsed.GetPath(focusInfo, supportingInfos)
	//
	// println(focusInfo.Uid)
	// for _, e := range supportingInfos {
	// 	println(e.Uid)
	// }
	_ = extraInfos
	_ = pth
	tm.LastFocus = focusInfo
	tm.LastSupportingInfos = supportingInfos
	tm.LastPath = pth
	return pth
	// pth.Print()
	//TODO: add section which runs through the available river and cascades to verify compliance UPDATE: not in this function anymore
	//TODO: if compliant, update model which the decision making module also uses. <- IMPORTANT UPDATE: not in this function anymore

}
