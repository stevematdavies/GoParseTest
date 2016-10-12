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

// Counter - Represents an entire Counter Object
type Counter struct {
	StartTime       string
	BaseID          string
	MeasurementID   string
	MeasurementType string
	Measurements    map[string]string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	counterMap := make(map[string][]string)
	data := getCountersData()
	counterMap["PmCounter"] = data
	fmt.Println(data)
}

func conStruct(rootElement Root) []Counter {

	pmSetupArray := rootElement.PMSetup
	var countersList []Counter

	for i := 0; i < len(pmSetupArray); i++ {
		counter := Counter{}
		pmSetup := pmSetupArray[i]
		measurements := pmSetup.MeasurementOutputResult.MeasurementList.Measurements
		measurementsMap := make(map[string]string)
		counter.StartTime = pmSetup.StartTime
		counter.BaseID = pmSetup.MeasurementOutputResult.MeasurementOutput.BaseID
		counter.MeasurementID = pmSetup.MeasurementOutputResult.MeasurementOutput.LocalMoID
		counter.MeasurementType = pmSetup.MeasurementOutputResult.MeasurementList.Type

		for j := 0; j < len(measurements); j++ {
			measure := measurements[j]
			name := measure.XMLName.Local
			content := measure.Content
			measurementsMap[name] = content
		}

		counter.Measurements = measurementsMap
		countersList = append(countersList, counter)
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
