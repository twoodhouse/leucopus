package timeKeeper

import (
	"github.com/twoodhouse/leucopus/decisionMaker"
	"github.com/twoodhouse/leucopus/info"
	"github.com/twoodhouse/leucopus/memory"
	"github.com/twoodhouse/leucopus/model"
	"github.com/twoodhouse/leucopus/testMaster"
)

type TimeKeeper struct {
	Mem *memory.Memory
	Tm  *testMaster.TestMaster
	Dm  *decisionMaker.DecisionMaker
	Mdl *model.Model
}

func New() *TimeKeeper {
	tm := testMaster.New()
	dm := decisionMaker.New()
	mdl := model.New()
	var entity = TimeKeeper{
		tm.Mem,
		tm,
		dm,
		mdl,
	}
	return &entity
}

func (tk *TimeKeeper) InitInfo(uid string) {
	nfo := info.New(uid)
	tk.Mem.SetRiver(nfo, []int{})
}

/*
Notes:
TODO list:
- Path algorithm
  - Currently switches forward to the next supporting info set at after a constant number of attempts
	- Currently the algorithm is taking (possibly) too long to reach higher numbers of I nodes
*/
func (tk *TimeKeeper) Begin() {
	for i := 0; i < 1; i++ { //TODO: increase this iteration value as appropriate for testing
		//*** Complete any decision maker actions here ***
		// assign action info values for this iteration

		//*** Complete any Memory actions here
		// apply action info values and input info values

		for j := 0; j < 15; j++ {
			//*** Complete any test master actions here ***
			// determine path to try (along with necessary focus and support)
			pth := tk.Tm.GetNextPath()
			pth.Print()
			println()
			_ = pth
			// try path against current cascades
			var fitGoodness float32
			fitGoodness = tk.Mem.ProcessPathAgainstCascades(pth)
			// update model
			tk.Mdl.CheckInPath(tk.Tm.LastFocus, tk.Tm.LastSupportingInfos, pth, fitGoodness)
			// update memory if appropriate

		}
	}
}
