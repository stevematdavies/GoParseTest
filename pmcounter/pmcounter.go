package pmcounter

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

// Root ..
type Root struct {
	XMLName xml.Name  `xml:"OMeS"`
	PMSetup []PMSetup `xml:"PMSetup"`
}

// PMSetup ...
type PMSetup struct {
	XMLName                 xml.Name                `xml:"PMSetup"`
	Interval                string                  `xml:"interval,attr"`
	StartTime               string                  `xml:"startTime,attr"`
	MeasurementOutputResult MeasurementOutputResult `xml:"PMMOResult"`
}

// MeasurementOutputResult ...
type MeasurementOutputResult struct {
	XMLName           xml.Name `xml:"PMMOResult"`
	MeasurementOutput MeasurementOutput
	MeasurementList   MeasurementList `xml:"NE-WBTS_1.0"`
}

// MeasurementOutput ...
type MeasurementOutput struct {
	XMLName   xml.Name `xml:"MO"`
	BaseID    string   `xml:"baseId"`
	LocalMoID string   `xml:"localMoid"`
}

// MeasurementList ...
type MeasurementList struct {
	XMLName      string        `xml:"NE-WBTS_1.0"`
	Type         string        `xml:"measurementType,attr"`
	Measurements []Measurement `xml:",any"`
}

// Measurement ...
type Measurement struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

type CounterMeasurement struct {
	CounterID string
	Value     string
}

// Counter - Represents an entire Counter Object
type Counter struct {
	CountTime       string
	DeviceID        string
	Device          string
	MeasurementType string
	Counters        []CounterMeasurement
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var kuhaID string

// GetPMCountersForDevice start here for FZM PM COounters
func GetPMCountersForDevice(kid string) []string {
	kuhaID = kid
	return getCountersData()
}

/*
* Temp list of required Measurement Types for GetCounters
* This will eventually be removed and the types
* will be filtered in the FZM before being parsed
 */
func typeInFilterList(mt string) bool {
	switch mt {
	case
		"LTE_Cell_Avail",
		"LTE_SINR",
		"LTE_RRC",
		"LTE_Cell_Throughput",
		"LTE_Cell_Load":
		return true
	}
	return false
}

func conStruct(rootElement Root) []Counter {

	pmSetupArray := rootElement.PMSetup
	var countersList []Counter

	for i := 0; i < len(pmSetupArray); i++ {
		counter := Counter{}
		pmSetup := pmSetupArray[i]
		measurements := pmSetup.MeasurementOutputResult.MeasurementList.Measurements
		counter.CountTime = pmSetup.StartTime
		counter.DeviceID = kuhaID
		counter.Device = pmSetup.MeasurementOutputResult.MeasurementOutput.LocalMoID
		counter.MeasurementType = pmSetup.MeasurementOutputResult.MeasurementList.Type
		measurementsArray := make([]CounterMeasurement, len(measurements))

		for j := 0; j < len(measurements); j++ {
			measurement := measurements[j]
			counterMeasurement := CounterMeasurement{}
			counterMeasurement.CounterID = measurement.XMLName.Local
			counterMeasurement.Value = measurement.Content
			measurementsArray[j] = counterMeasurement

		}

		counter.Counters = measurementsArray

		if typeInFilterList(counter.MeasurementType) {
			countersList = append(countersList, counter)
			fmt.Println(counter.MeasurementType)
		}
	}
	return countersList
}

func xmlify(byts []byte) []Counter {
	root := Root{}
	err := xml.Unmarshal(byts, &root)
	check(err)
	return conStruct(root)
}

func getCountersData() []string {
	byts, err := ioutil.ReadFile(path())
	check(err)
	return jsonFromXML(xmlify(byts))
}

func path() string {
	pmXMLFilePath := "/tmp/stats/pm.counter.xml"
	matches, err := filepath.Glob("/tmp/stats/*.raw.xz")
	check(err)
	if len(matches) > 0 {
		data, err := ioutil.ReadFile(matches[0])
		check(err)
		readable, err := xz.NewReader(bytes.NewReader(data), 0)
		check(err)
		var byts bytes.Buffer
		_, err = byts.ReadFrom(readable)
		check(err)
		err = ioutil.WriteFile(pmXMLFilePath, byts.Bytes(), 0644)
		check(err)
		return pmXMLFilePath
	}
	return ""
}

func jsonFromXML(countersArray []Counter) []string {
	var JSONCounterArray []string
	for i := 0; i < len(countersArray); i++ {
		counter := countersArray[i]
		jsn, err := json.Marshal(counter)
		check(err)
		JSONCounterArray = append(JSONCounterArray, string(jsn)+",")
	}
	return stripArrayTrailingComma(JSONCounterArray)
}

func stripArrayTrailingComma(arr []string) []string {
	index := len(arr) - 1
	arr[index] = strings.TrimSuffix(arr[index], ",")
	return arr
}
