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
	extraInfos := t.SupportingInfosModuleUsed.GetExtraInfos(focusInfo) //TODO: this has something to do with cascades?
	//TODO: extraInfos need to be added to the cascade if the goop fails
	pth := t.PathModuleUsed.GetPath(focusInfo, supportingInfos)

	_ = extraInfos
	_ = pth
	// pth.Print()
	//TODO: add section which runs through the available river and cascades to verify compliance
	//TODO: if compliant, update model which the decision making module also uses. <- IMPORTANT

}
