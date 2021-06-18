package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"io/ioutil"
	"strconv"
	"time"
)

// 	https://stats.kmd.io/api/source/mined/?name=slyris_EU&from=timestamp&to=timestamp
//  https://stats.kmd.io/api/source/mined/?name=slyris_EU&min_blocktime=1623861228&max_blocktime=1623893930

const (
	layoutISO = "2006-01-02"
	layoutUS  = "January 2, 2006"
)

type Config struct {
	NotaryID []string `json:"notary_id"`
	Range    []string `json:"range"`
}

type StatsKmdMining struct {
	Count    int         `json:"count"`
	Next     interface{} `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		BlockHeight   int       `json:"block_height"`
		BlockTime     int       `json:"block_time"`
		BlockDatetime time.Time `json:"block_datetime"`
		Value         string    `json:"value"`
		Address       string    `json:"address"`
		Name          string    `json:"name"`
		Txid          string    `json:"txid"`
		Season        string    `json:"season"`
	} `json:"results"`
}

var gConfig Config
var gStatsMining StatsKmdMining

func parseConfig() {
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &gConfig)
	fmt.Printf("%v\n", gConfig)
}

func main() {
	parseConfig()
	for _, s := range gConfig.NotaryID {
		fmt.Printf("Calculating taxe revenue of %s\n", s)
		for _, curRange := range gConfig.Range {
			switch curRange {
			case "month":
				calculateTaxeLastMonth(s)
			case "year":
				calculateTaxeLastYear(s)
			}
		}
	}
}

func fillMiningStats(notaryNodeId string, from time.Time, to time.Time) {
	//fmt.Printf("%s\n", from.Format(layoutISO))
	//fmt.Printf("%s\n", to.Format(layoutISO))
	url := "https://stats.kmd.io/api/source/mined/?name=" + notaryNodeId + "&min_blocktime=" + strconv.FormatInt(from.Unix(), 10) + "&max_blocktime=" + strconv.FormatInt(to.Unix(), 10)
	fmt.Printf("Processing: %s\n", url)
}

func calculateTaxeLastYear(notaryNodeId string) {
	fmt.Printf("Calculating tax for last year: [%s]\n", notaryNodeId)
	beginningOfYear := now.BeginningOfYear()
	endOfYear := now.EndOfYear()
	fmt.Printf("timestamp begin of the year: %v %d\n", beginningOfYear, beginningOfYear.Unix())
	fmt.Printf("timestamp end of the year: %v %d\n", endOfYear, endOfYear.Unix())
	fillMiningStats(notaryNodeId, beginningOfYear, endOfYear)
}

func calculateTaxeLastMonth(notaryNodeId string) {
	fmt.Printf("Calculating tax for last month: [%s]\n", notaryNodeId)
	beginningOfMonth := now.BeginningOfMonth()
	endOfMonth := now.EndOfMonth()
	fmt.Printf("timestamp begin of the month: %v %d\n", beginningOfMonth, beginningOfMonth.Unix())
	fmt.Printf("timestamp end of the month: %v %d\n", endOfMonth, endOfMonth.Unix())
	fillMiningStats(notaryNodeId, beginningOfMonth, endOfMonth)
}
