package info

//************
type Info struct { //each info is correlated directly with an input to the system and given a uid
	Fit float32
	Uid string
}

func New(uid string) *Info {
	var entity = Info{
		0, //TODO: fit
		uid,
	}
	return &entity
}

func (info *Info) GrabCurrentValue() int {
	//TODO: implement this based off another class or system
	return 0
}

//
// //*************
// type River struct { //the river is the set of information which is fed to the pather for analysis. This is the full set of remembered info.
// 	PrimaryFlow *Flow
// 	Flows       []*Flow
// 	Cascades    []*Cascade
// }
//
// func (river *River) NewRiver() *River {
// 	var entity = River{
// 		0, //TODO: fit
// 		uid,
// 	}
// 	return &entity
// }
//
// //*************
// type Flow struct { //each flow is a separate set of information which the tables must be tested against
// 	Streams []*Stream
// }
//
// func (flow *Flow) NewFlow() *Flow {
// 	var entity = Flow{
// 		0, //TODO: fit
// 		uid,
// 	}
// 	return &entity
// }
//
// //*************
// type Cascade struct { //each cascade is a limited set of history across limited streams captured due to nonconformity
// 	Streams []*Stream
// }
//
// func (cascade *Cascade) NewCascade() *Cascade {
// 	var entity = Cascade{
// 		0, //TODO: fit
// 		uid,
// 	}
// 	return &entity
// }
//
// //*************
// type Stream struct { //each stream is a history of outputs related to a single Info
// 	Info    *Info
// 	History []int
// }
//
// func (stream *Stream) NewSteam(info *Info) *Steam {
// 	var entity = Steam{
// 		info,
// 		[]int{},
// 	}
// 	return &entity
// }
