package pmcounter

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"
	"path/filepath"
	"strings"

	"xi2.org/x/xz"
)

// Root ..
type Root struct {
	XMLName        xml.Name        `xml:"OMeS"`
	ParentElements []ParentElement `xml:"PMSetup"`
}

// ParentElement ...
type ParentElement struct {
	XMLName             xml.Name            `xml:"PMSetup"`
	ParentElementResult ParentElementResult `xml:"PMMOResult"`
}

// ParentElementResult ...
type ParentElementResult struct {
	XMLName           xml.Name          `xml:"PMMOResult"`
	DeviceCounterList DeviceCounterList `xml:"NE-WBTS_1.0"`
}

// DeviceCounterList ...
type DeviceCounterList struct {
	DeviceCounterArray []DeviceCounter `xml:",any"`
}

// DeviceCounter ...
type DeviceCounter struct {
	XMLName xml.Name
	Content string `xml:",innerxml"`
}

// DeviceCounterJSON ...
type DeviceCounterJSON struct {
	Name    string
	Content string
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// GetPMCounters start here for FZM PM COounters
func GetPMCounters() []string {
	return getCountersData()
}

func conStruct(rootElement Root) []DeviceCounterJSON {
	parentElements := rootElement.ParentElements
	var deviceCountersList []DeviceCounterJSON
	for i := 0; i < len(parentElements); i++ {
		pmCounters := parentElements[i].ParentElementResult.DeviceCounterList.DeviceCounterArray
		for j := 0; j < len(pmCounters); j++ {
			dcj := DeviceCounterJSON{}
			dcj.Name = pmCounters[j].XMLName.Local
			dcj.Content = pmCounters[j].Content
			deviceCountersList = append(deviceCountersList, dcj)
		}
	}
	return deviceCountersList
}

func xmlify(byts []byte) []DeviceCounterJSON {
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

func jsonFromXML(countersArray []DeviceCounterJSON) []string {
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
