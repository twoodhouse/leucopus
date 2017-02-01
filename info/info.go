package info

type Info struct { //each info is correlated directly with an input to the system and given a uid
	Fit float32
	Uid string
	Url string
}

func New(uid string, url string) *Info {
	var entity = Info{
		0, //TODO: fit
		uid,
		url,
	}
	return &entity
}

// func (info *Info) Copy() {
// 	var newInfo *Info
// 	New()
// }

func (info *Info) GrabCurrentValue() int {
	//TODO: implement this based off another class or system
	return 0
}
