package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"xi2.org/x/xz"
)

type Root struct {
	XMLName xml.Name  `xml:"OMeS"`
	PMSetup []PMSetup `xml:"PMSetup"`
}

type PMSetup struct {
	XMLName    xml.Name   `xml:"PMSetup"`
	Interval   string     `xml:"interval,attr"`
	StartTime  string     `xml:"startTime,attr"`
	PMMOResult PMMOResult `xml:"PMMOResult"`
}

type PMMOResult struct {
	XMLName xml.Name `xml:"PMMOResult"`
	MO      MO
	NEWBTS  NEWBTS `xml:"NE-WBTS_1.0"`
}

type NEWBTS struct {
	XMLName         string          `xml:"NE-WBTS_1.0"`
	MeasurementType string          `xml:"measurementType,attr"`
	Counters        []DeviceCounter `xml:",any"`
}

type DeviceCounter struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type MO struct {
	XMLName   xml.Name `xml:"MO"`
	BaseID    string   `xml:"baseId"`
	LocalMoID string   `xml:"localMoid"`
}

type Counter struct {
	StartTime       string
	BaseID          string
	MeasurementType string
	WbitCounts      map[string]string
}

/* ************************************************************************************************ */

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	data := getCountersData()
	fmt.Println(data)
}

func conStruct(r Root) []Counter {

	pmSetups := r.PMSetup

	var ctrscol []Counter

	for i := 0; i < len(pmSetups); i++ {
		ctr := Counter{}
		ps := pmSetups[i]
		ctrArr := ps.PMMOResult.NEWBTS.Counters
		ctrMap := make(map[string]string)
		ctr.StartTime = ps.StartTime
		ctr.BaseID = ps.PMMOResult.MO.BaseID
		ctr.MeasurementType = ps.PMMOResult.NEWBTS.MeasurementType

		for j := 0; j < len(ctrArr); j++ {
			dctr := ctrArr[j]
			dcName := dctr.XMLName.Local
			dcContent := dctr.Content
			ctrMap[dcName] = dcContent
		}

		ctr.WbitCounts = ctrMap
		ctrscol = append(ctrscol, ctr)
	}
	return ctrscol
}

func xmlify(f []byte) []Counter {
	root := Root{}
	err := xml.Unmarshal(f, &root)
	check(err)
	return conStruct(root)
}

func getCountersData() []string {
	byts, err := ioutil.ReadFile(path())
	check(err)
	return jsonFromXML(xmlify(byts))
}

func path() string {
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

func jsonFromXML(countersArray []Counter) []string {
	var JSONCounterArray []string
	for i := 0; i < len(countersArray); i++ {
		c := countersArray[i]
		jsn, err := json.Marshal(c)
		check(err)
		JSONCounterArray = append(JSONCounterArray, string(jsn)+",")
	}
	return stripArrayTrailingComma(JSONCounterArray)
}

func stripArrayTrailingComma(arr []string) []string {
	i := len(arr) - 1
	arr[i] = strings.TrimSuffix(arr[i], ",")
	return arr
}
