package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"xi2.org/x/xz"
)

type OMeS struct {
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

func check(err error) {
	if err != nil {
		panic(err)
	}
}

type Counter struct {
	StartTime       string
	BaseID          string
	MeasurementType string
	WbitCounts      map[string]string
}

func main() {
	getData()

}

type DeviceCounterTag struct {
	Name  string
	Value string
}

func construct(r OMeS) {
	pmSetups := r.PMSetup
	for i := 0; i < len(pmSetups); i++ {
		ctr := Counter{}
		ps := pmSetups[i]
		ctrArr := ps.PMMOResult.NEWBTS.Counters
		ctr.StartTime = ps.StartTime
		ctr.BaseID = ps.PMMOResult.MO.BaseID
		ctr.MeasurementType = ps.PMMOResult.NEWBTS.MeasurementType

		for j := 0; j < len(ctrArr); j++ {
			dctr := ctrArr[j]
			dcName := dctr.XMLName.Local
			dcContent := dctr.Content

			fmt.Printf("%s  :  %s\n", dcName, dcContent)

		}

	}

}

func parseData(f []byte) {
	root := OMeS{}
	err := xml.Unmarshal(f, &root)
	check(err)
	fmt.Println(root)
	construct(root)
}

func getData() {
	xml, err := ioutil.ReadFile(getCounterFile())
	check(err)
	parseData(xml)
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
