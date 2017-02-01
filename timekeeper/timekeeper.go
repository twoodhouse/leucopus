package timeKeeper

import (
	"github.com/twoodhouse/leucopus/decisionMaker"
	"github.com/twoodhouse/leucopus/ioManager"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/model"
	"github.com/twoodhouse/leucopus/testMaster"
)

type TimeKeeper struct {
	Mem *memory.Memory
	Tm  *testMaster.TestMaster
	Dm  *decisionMaker.DecisionMaker
	Mdl *model.Model
	Iom *ioManager.IoManager
}

func New(actionUrls []string, observeUrls []string, masterUrl string) *TimeKeeper {
	tm := testMaster.New()
	mdl := model.New()
	iom := ioManager.New(actionUrls, observeUrls, masterUrl)
	dm := decisionMaker.New(mdl, iom)
	var tk = TimeKeeper{
		tm.Mem,
		tm,
		dm,
		mdl,
		iom,
	}
	for _, actionUrl := range actionUrls {
		tk.InitAction(actionUrl)
	}
	for _, observeUrl := range observeUrls {
		tk.InitInfo(observeUrl)
	}
	tk.Mem.SetRiver(iom.MasterInfo, []int{})
	tk.Tm.FocusInfoModuleUsed.PullRiverInfos()
	tk.Tm.SupportingInfosModuleUsed.PullRiverInfos()
	return &tk
}

//TODO: functionality not verified
func (tk *TimeKeeper) InitInfo(url string) {
	nfo := tk.Iom.MakeInfo(url)
	tk.Mem.SetRiver(nfo, []int{})
}

//TODO: functionality not verified
func (tk *TimeKeeper) InitAction(url string) {
	nfo := tk.Iom.MakeAction(url)
	tk.Mem.SetRiver(nfo, []int{})
}

/*
Notes:
TODO list:
- Path algorithm
  - Currently switches forward to the next supporting info set at after a constant number of attempts
	- Currently the algorithm is taking (possibly) too long to reach higher numbers of I nodes
		- This could be corrected by saying a certain number of attempts must be made at a lower level, then the next is tried (recursive)
	- Method for choosing supporting I nodes may need revised
	- Investigate ExitINode issues (MAJOR)
	- Pushing buttons and checking them in different orders could cause problems between runs
*/
func (tk *TimeKeeper) Begin() {
	tk.Mem.PrintRiver()

	for i := 0; i < 2; i++ { //TODO: increase this iteration value as appropriate for testing
		//*** Complete any decision maker actions here ***
		// assign action info values for this iteration
		actionRowMap := tk.Dm.GetIncomingActionRowMap()

		//*** Complete any Memory actions here
		// get world info values from IOManager
		fullRowMap := tk.Iom.GetIncomingInfoAndMasterRowMap()
		// take actions supplied by decision maker
		tk.Iom.ProcessActionRowMap(actionRowMap)
		// apply action info values and input info values
		for k, v := range actionRowMap {
			fullRowMap[k] = v
		}
		tk.Mem.MagicRiverInput(fullRowMap)

		for j := 0; j < 3; j++ {
			//*** Complete any test master actions here ***
			// determine path to try (along with necessary focus and support)
			pth := tk.Tm.GetNextPath()
			// try path against current cascades
			var unweightedGoodness float32
			unweightedGoodness = tk.Mem.ProcessPathAgainstCascades(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos, pth) //TODO: this needs verified
			if unweightedGoodness > tk.Mdl.GetFitGoodness(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos).UnweightedGoodness {
				// update model if appropriate
				tk.Mdl.CheckInPath(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos, pth, unweightedGoodness)
				// update memory if appropriate
				// println()
				// println(tk.Tm.LastFocus.Uid)
				// println(unweightedGoodness)
				// println(tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus])
				if unweightedGoodness > tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] || tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] == float32(0) { //TODO: consider changing fitGoodness so that it works from 1 to 0, rather than 0 to 1
					tk.Mem.Paths[tk.Tm.LastFocus] = pth
					tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] = unweightedGoodness
					tk.Mem.SupportingInfos[tk.Tm.LastFocus] = tk.Tm.LastSupportingInfos
					//TODO: LOW PRIORITY - add depth if appropriate later on
					//TODO: implement cascade functionality for adding common infos to track
				}
			}
		}
	}
	tk.Mem.PrintPaths()
	tk.Mem.PrintGoodness()
	tk.Mdl.Print()
	tk.Mem.PrintRiver()
	tk.Mem.PrintNumCascades() //TODO: investigate why no cascades are being created
}
