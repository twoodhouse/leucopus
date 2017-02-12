package timeKeeper

import (
	"math"

	"github.com/twoodhouse/leucopus/decisionMaker"
	"github.com/twoodhouse/leucopus/ioManager"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/model"
	"github.com/twoodhouse/leucopus/testMaster"
	"github.com/twoodhouse/leucopus/utility"
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
	- TODO next: make the value itself not available for checking? Or the whole value row? Or a special setting for doing current checks with other sources?
	- Pushing buttons and checking them in different orders could cause problems between runs
*/
func (tk *TimeKeeper) Begin() {
	for i := 0; i < 100; i++ { //TODO: increase this iteration value as appropriate for testing
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
		//check in each optimum path to the model since there were likely changes
		for nfo, pth := range tk.Mem.Paths {
			tk.Mdl.CheckInPath(nfo, tk.Mem.SupportingInfos[nfo], pth, tk.Mem.InfoFitGoodnesses[nfo])
		}
		//*******************************************************************************************************************
		for j := 0; j < 4; j++ { //TODO: modify this section so that getting a path is just an alternative to making a swap to an old model.
			//iterate over all the infos and see if any existing models can be swapped in.
			foundModelPath := false
			for nfo, _ := range tk.Mem.River {
				if tk.Mem.InfoFitGoodnesses[nfo] < 1 && nfo.Uid == "M0" { //TODO: remove portion that limits this analysis to the master info
					pth, supportingInfos := tk.Mdl.FindPathBetterThanGoodness(nfo, tk.Mem.InfoFitGoodnesses[nfo])
					if pth != nil {
						unweightedGoodness := tk.Mem.ProcessPathAgainstCascades(nfo, supportingInfos, pth)
						foundModelPath = true
						// best model fit for each supporting info combo needs compared against the current set of cascades also
						modelUnweightedGoodness := tk.Mem.ProcessPathAgainstCascades(nfo, supportingInfos, tk.Mdl.Library[nfo][utility.GetUidFromInfos(supportingInfos)])
						tk.Mdl.CheckInPath(nfo, supportingInfos, tk.Mdl.Library[nfo][utility.GetUidFromInfos(supportingInfos)], modelUnweightedGoodness)
						//look for path updates
						if unweightedGoodness > tk.Mem.InfoFitGoodnesses[nfo] || tk.Mem.InfoFitGoodnesses[nfo] == float32(0) {
							tk.Mem.Paths[nfo] = pth
							tk.Mem.RiverBalance[nfo] = make([]int, int(math.Exp2(float64(len(tk.Mem.Paths[nfo].MiddleLinks[len(tk.Mem.Paths[nfo].MiddleLinks)-1].Inputs)))))
							tk.Mem.InfoFitGoodnesses[nfo] = unweightedGoodness
							tk.Mem.SupportingInfos[nfo] = supportingInfos
						}
					}
				}
			}
			if !foundModelPath {
				// determine path to try (along with necessary focus and support)
				pth := tk.Tm.GetNextPath()

				if tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] < 1 && tk.Tm.LastFocus.Uid == "M0" { //TODO: remove portion that limits this analysis to the master info
					pth.Print()
					println(tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus])
					// try path against current cascades
					var unweightedGoodness float32
					unweightedGoodness = tk.Mem.ProcessPathAgainstCascades(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos, pth)
					if tk.Tm.LastFocus.Uid == "M0" { //for test only
						println(unweightedGoodness)
					}
					if unweightedGoodness > tk.Mdl.GetFitGoodness(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos).UnweightedGoodness {
						// update model if appropriate
						tk.Mdl.CheckInPath(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos, pth, unweightedGoodness)
					}
					//look for path updates
					if unweightedGoodness > tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] || tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] == float32(0) {
						tk.Mem.Paths[tk.Tm.LastFocus] = pth
						tk.Mem.RiverBalance[tk.Tm.LastFocus] = make([]int, int(math.Exp2(float64(len(tk.Mem.Paths[tk.Tm.LastFocus].MiddleLinks[len(tk.Mem.Paths[tk.Tm.LastFocus].MiddleLinks)-1].Inputs)))))
						tk.Mem.InfoFitGoodnesses[tk.Tm.LastFocus] = unweightedGoodness
						tk.Mem.SupportingInfos[tk.Tm.LastFocus] = tk.Tm.LastSupportingInfos
						//TODO: LOW PRIORITY - add depth if appropriate later on
						//TODO: implement cascade functionality for adding common infos to track
					}
				}
			}
		}
		//*******************************************************************************************************************

	}
	tk.Mem.PrintPaths()
	tk.Mem.PrintGoodness()
	tk.Mdl.Print()
	tk.Mem.PrintRiver()
	tk.Mem.PrintNumCascades() //TODO: investigate why no cascades are being created
}
