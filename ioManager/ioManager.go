package ioManager

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/twoodhouse/leucopus/info"
)

type IoManager struct {
	ActionInfos  []*info.Info
	ObserveInfos []*info.Info
	MasterInfo   *info.Info
	UidCounter   int
}

func New(actionUrls []string, observeUrls []string, masterUrl string) *IoManager {
	var entity = IoManager{
		[]*info.Info{},
		[]*info.Info{},
		info.New("M0", masterUrl),
		1,
	}
	// for _, actionUrl := range actionUrls {
	// 	entity.MakeAction(actionUrl)
	// }
	// for _, observeUrl := range observeUrls {
	// 	entity.MakeInfo(observeUrl)
	// }
	return &entity
}

//TODO: functionality not verified
func (iom *IoManager) MakeAction(url string) *info.Info {
	newInfo := info.New("A"+strconv.Itoa(iom.UidCounter), url)
	iom.ActionInfos = append(iom.ActionInfos, newInfo)
	iom.UidCounter = iom.UidCounter + 1
	return newInfo
}

//TODO: functionality not verified
func (iom *IoManager) MakeInfo(url string) *info.Info {
	newInfo := info.New("B"+strconv.Itoa(iom.UidCounter), url)
	iom.ObserveInfos = append(iom.ObserveInfos, newInfo)
	iom.UidCounter = iom.UidCounter + 1
	println("making info")
	return newInfo
}

func (iom *IoManager) GetIncomingInfoAndMasterRowMap() map[*info.Info]int {
	infoRowMap := make(map[*info.Info]int)
	resp, _ := http.Get(iom.MasterInfo.Url)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	bResp, _ := strconv.ParseBool(string(body))
	iResp := Btoi(bResp)
	infoRowMap[iom.MasterInfo] = iResp
	for _, observeInfo := range iom.ObserveInfos {
		resp, _ := http.Get(observeInfo.Url)
		defer resp.Body.Close()
		body, _ := ioutil.ReadAll(resp.Body)
		bResp, _ := strconv.ParseBool(string(body))
		iResp := Btoi(bResp)
		infoRowMap[observeInfo] = iResp
	}
	return infoRowMap
}

func (iom *IoManager) ProcessActionRowMap(actionRowMap map[*info.Info]int) {
	for _, actionInfo := range iom.ActionInfos {
		if actionRowMap[actionInfo] == 1 {
			http.Get(actionInfo.Url)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func Btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
