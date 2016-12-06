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
}

type FocusInfoModule interface {
	GetFocusInfo() *info.Info
}

type SupportingInfosModule interface {
	GetSupportingInfos(*info.Info) []*info.Info
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
	}
	return &entity
}

func (t *TestMaster) GoopNextPath() {
	focusInfo := t.FocusInfoModuleUsed.GetFocusInfo()
	supportingInfos := t.SupportingInfosModuleUsed.GetSupportingInfos(focusInfo)
	extraInfos := t.SupportingInfosModuleUsed.GetExtraInfos(focusInfo)
	pth := t.PathModuleUsed.GetPath(focusInfo, supportingInfos)
	// pth.Print()
	//TODO: extraInfos need to be added to the cascade if the goop fails
	_ = extraInfos
	_ = supportingInfos
	_ = pth
}
