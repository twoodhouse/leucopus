package model

import (
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/pather"
	"github.com/twoodhouse/leucopus/utility"
)

type Model struct {
	//(info)(supportingInfos) -> Path... Normally, I would only have to gather paths by info, but varied supporting infos allows full modeling
	Library     map[*info.Info]map[string]*pather.Path
	FitGoodness map[*info.Info]map[string]*FitGoodness
	/* problem right now:
	the data captured in a cascade is not ALL information in the river, only certain infos.
	I suppose there is no real harm in capturing the whole moment. It just might mean I have to delete
	older cascades sooner.
	Changes due to this modification:
	1. All of the River is stored in a cascade when it is opened
	2. Cascades are only indexed according to focusInfo, rather than focus and supporting
	*/
}

type FitGoodness struct {
	UnweightedGoodness float32
	Attempts           int
}

func New() *Model {
	var entity = Model{
		make(map[*info.Info]map[string]*pather.Path),
		make(map[*info.Info]map[string]*FitGoodness),
	}
	return &entity
}

func NewFitGoodness() *FitGoodness {
	var entity = FitGoodness{
		0,
		0,
	}
	return &entity
}

func (mdl *Model) CheckInPath(focusInfo *info.Info, supportingInfos []*info.Info, pth *pather.Path, unweightedGoodness float32) {
	if _, ok := mdl.Library[focusInfo]; !ok {
		mdl.Library[focusInfo] = make(map[string]*pather.Path)
	}
	mdl.Library[focusInfo][utility.GetUidFromInfos(supportingInfos)] = pth

	if _, ok := mdl.FitGoodness[focusInfo]; !ok {
		mdl.FitGoodness[focusInfo] = make(map[string]*FitGoodness)
	}

	if _, ok := mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)]; !ok {
		fg := NewFitGoodness()
		fg.UnweightedGoodness = unweightedGoodness
		fg.Attempts = 1
		mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)] = fg
	} else {
		fgSelection := mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)]
		fgSelection.UnweightedGoodness = unweightedGoodness
		fgSelection.Attempts = fgSelection.Attempts + 1
	}
}

func (mdl *Model) GetFitGoodness(focusInfo *info.Info, supportingInfos []*info.Info) *FitGoodness {
	if _, ok := mdl.FitGoodness[focusInfo]; !ok {
		mdl.FitGoodness[focusInfo] = make(map[string]*FitGoodness)
	}
	if _, ok := mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)]; !ok {
		fg := NewFitGoodness()
		fg.UnweightedGoodness = 0
		fg.Attempts = 0
		mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)] = fg
	}
	return mdl.FitGoodness[focusInfo][utility.GetUidFromInfos(supportingInfos)]
}

func (mdl *Model) Print() {
	println()
	println("PRINTING MODEL FIT GOODNESS")
	for nfo, _ := range mdl.FitGoodness {
		println(nfo.Uid)
		print("{")
		for sNfos, fg := range mdl.FitGoodness[nfo] {
			print(sNfos)
			print(": [")
			print(fg.Attempts)
			print("]")
			print(fg.UnweightedGoodness)
			print(", ")
		}
		println("}")
	}
}
