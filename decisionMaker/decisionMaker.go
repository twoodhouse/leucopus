package decisionMaker

import (
	"math/rand"

	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/ioManager"
	"github.com/twoodhouse/leucopus/model"
)

type DecisionMaker struct {
	Mdl *model.Model
	Iom *ioManager.IoManager
}

func New(mdl *model.Model, iom *ioManager.IoManager) *DecisionMaker {
	var entity = DecisionMaker{
		mdl,
		iom,
	}
	return &entity
}

func (dm *DecisionMaker) GetIncomingActionRowMap() map[*info.Info]int {
	infoRowMap := make(map[*info.Info]int)
	//initial implementation to allow for other work first: Random flailing. TODO: change this into rational decision making
	for _, actionInfo := range dm.Iom.ActionInfos {
		infoRowMap[actionInfo] = rand.Intn(2)
		// infoRowMap[actionInfo] = 1
	}
	return infoRowMap
}
