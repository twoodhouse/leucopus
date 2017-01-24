package model

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/utility"
)

type Model struct {
	//(info)(supportingInfos) -> Path... Normally, I would only have to gather paths by info, but varied supporting infos allows full modeling
	Library     map[*info.Info]map[string]*pather.Path
	FitGoodness map[*info.Info]map[string]float32
	/* problem right now:
	the data captured in a cascade is not ALL information in the river, only certain infos.
	I suppose there is no real harm in capturing the whole moment. It just might mean I have to delete
	older cascades sooner.
	Changes due to this modification:
	1. All of the River is stored in a cascade when it is opened
	2. Cascades are only indexed according to focusInfo, rather than focus and supporting
	*/
}

func New() *Model {
	var entity = Model{
		make(map[*info.Info]map[string]*pather.Path),
		make(map[*info.Info]map[string]float32),
	}
	return &entity
}

func (mdl *Model) CheckInPath(focusInfo *info.Info, supportingInfos []*info.Info, pth *pather.Path, fitGoodness float32) {
	if _, ok := mdl.Library[focusInfo]; !ok {
		mdl.Library[focusInfo] = make(map[string]*pather.Path)
	}
	mdl.Library[focusInfo][utility.GetUidFromInfos(supportingInfos)] = pth

	if _, ok := mdl.FitGoodness[focusInfo]; !ok {
		mdl.FitGoodness[focusInfo] = make(map[string]float32)
	}
	mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)] = fitGoodness
}
