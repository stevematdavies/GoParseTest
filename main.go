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
	BbaseID   string `json:"baseId"`
	LocalMoID string `json:"localMoid"`
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	getJSONData()

}

func parseJSON(j string) {
	root := Root{}
	err := json.Unmarshal([]byte(j), &root)
	check(err)
	fmt.Printf("%v", root.OMeS.PMSetup[0].PMMOResult.NEWBTS)
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
