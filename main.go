package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"

	xj "github.com/basgys/goxml2json"
	"xi2.org/x/xz"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	jsonData, err := getJSONData()
	check(err)
	parseJSON(jsonData)
}

func parseJSON(j *bytes.Buffer) {
	// TODO lets get tagging!
}

func getJSONData() (*bytes.Buffer, error) {
	file, err := os.Open(getCounterFile())
	check(err)
	defer file.Close()
	return xj.Convert(file)
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
