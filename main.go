package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	xj "github.com/basgys/goxml2json"
	"xi2.org/x/xz"
)

type Root struct {
	OMeS OMeS
}

type OMeS struct {
	PMSetup []PMSetup
}

type PMSetup struct {
	Interval   string     `json:"-interval"`
	StartTime  string     `json:"-startTime"`
	PMMOResult PMMOResult `json:"PMMOResult"`
}

type PMMOResult struct {
	MO     MO
	NEWBTS interface{} `json:"NE-WBTS_1.0"`
}

type MO struct {
	BaseID    string `json:"baseId"`
	LocalMoID string `json:"localMoid"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Counter struct {
	StartTime       string
	BaseID          string
	MeasurementType string
	WbitCounts      map[string]interface{}
}

func main() {
	getJSONData()

}

func initCounter(r Root) {
	head := r.OMeS.PMSetup[0]
	//for i := 0; i < len(head); i++ {
	c := Counter{}
	h := head
	mesT := "-measurementType"
	c.StartTime = "Start Time: " + h.StartTime
	c.BaseID = "BaseID: " + h.PMMOResult.MO.BaseID
	ctrMap := h.PMMOResult.NEWBTS.(map[string]interface{})
	var mt = ctrMap[mesT].(string)
	c.MeasurementType = "Measurement Type: " + mt
	c.WbitCounts = make(map[string]interface{})
	for k, v := range ctrMap {
		if k != mesT {
			c.WbitCounts[k] = v.(string)
		}
	}

	fmt.Println(c)

}

func parseJSON(j string) {
	root := Root{}
	err := json.Unmarshal([]byte(j), &root)
	check(err)
	initCounter(root)
}

func getJSONData() {
	file, err := os.Open(getCounterFile())
	check(err)
	defer file.Close()
	json, err := xj.Convert(file)
	check(err)
	parseJSON(json.String())
}

func getCounterFile() string {
	pmFile := "/tmp/stats/pm.counter.xml"
	m, err := filepath.Glob("/tmp/stats/*.raw.xz")
	check(err)

	if len(m) > 0 {
		data, err := ioutil.ReadFile(m[0])
		check(err)
		rdr, err := xz.NewReader(bytes.NewReader(data), 0)
		check(err)
		var outBfr bytes.Buffer
		_, err = outBfr.ReadFrom(rdr)
		check(err)
		err = ioutil.WriteFile(pmFile, outBfr.Bytes(), 0644)
		check(err)
		return pmFile
	}
	return ""
}
