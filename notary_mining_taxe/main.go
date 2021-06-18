package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/kpango/glg"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"
)

// 	https://stats.kmd.io/api/source/mined/?name=slyris_EU&from=timestamp&to=timestamp
//  https://stats.kmd.io/api/source/mined/?name=slyris_EU&min_blocktime=1623861228&max_blocktime=1623893930
//  https://komodostats.com/api/notary/mined.json?nodename=slyris_EU&start=2021-01-01&end=2021-12-31
// 	https://komodostats.com/api/notary/mined.json?nodename=slyris_EU&start=2021-06-01&end=2021-06-30

const (
	layoutISO        = "2006-01-02"
	layoutUS         = "January 2, 2006"
	layoutCoingecko  = "02-01-2006"
	layoutOutCsv     = "02 Jan 06 15:04 CET"
	komodoStatsUri   = "https://komodostats.com/api/notary/mined.json?nodename="
	smkStatsUri      = "https://stats.kmd.io/api/source/mined/?name="
	coingeckoUri     = "https://api.coingecko.com/api/v3/coins/komodo/history?date="
	explorerTxUri    = "https://komodod.com/t/"
	explorerBlockUri = "https://komodod.com/b/"
)

type TotalUSD struct {
	time.Time
}

type CoingeckoHistoryResponse struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Image  struct {
		Thumb string `json:"thumb"`
		Small string `json:"small"`
	} `json:"image"`
	MarketData struct {
		CurrentPrice struct {
			Aed  float64 `json:"aed"`
			Ars  float64 `json:"ars"`
			Aud  float64 `json:"aud"`
			Bch  float64 `json:"bch"`
			Bdt  float64 `json:"bdt"`
			Bhd  float64 `json:"bhd"`
			Bmd  float64 `json:"bmd"`
			Bnb  float64 `json:"bnb"`
			Brl  float64 `json:"brl"`
			Btc  float64 `json:"btc"`
			Cad  float64 `json:"cad"`
			Chf  float64 `json:"chf"`
			Clp  float64 `json:"clp"`
			Cny  float64 `json:"cny"`
			Czk  float64 `json:"czk"`
			Dkk  float64 `json:"dkk"`
			Dot  float64 `json:"dot"`
			Eos  float64 `json:"eos"`
			Eth  float64 `json:"eth"`
			Eur  float64 `json:"eur"`
			Gbp  float64 `json:"gbp"`
			Hkd  float64 `json:"hkd"`
			Huf  float64 `json:"huf"`
			Idr  float64 `json:"idr"`
			Ils  float64 `json:"ils"`
			Inr  float64 `json:"inr"`
			Jpy  float64 `json:"jpy"`
			Krw  float64 `json:"krw"`
			Kwd  float64 `json:"kwd"`
			Lkr  float64 `json:"lkr"`
			Ltc  float64 `json:"ltc"`
			Mmk  float64 `json:"mmk"`
			Mxn  float64 `json:"mxn"`
			Myr  float64 `json:"myr"`
			Ngn  float64 `json:"ngn"`
			Nok  float64 `json:"nok"`
			Nzd  float64 `json:"nzd"`
			Php  float64 `json:"php"`
			Pkr  float64 `json:"pkr"`
			Pln  float64 `json:"pln"`
			Rub  float64 `json:"rub"`
			Sar  float64 `json:"sar"`
			Sek  float64 `json:"sek"`
			Sgd  float64 `json:"sgd"`
			Thb  float64 `json:"thb"`
			Try  float64 `json:"try"`
			Twd  float64 `json:"twd"`
			Uah  float64 `json:"uah"`
			Usd  float64 `json:"usd"`
			Vef  float64 `json:"vef"`
			Vnd  float64 `json:"vnd"`
			Xag  float64 `json:"xag"`
			Xau  float64 `json:"xau"`
			Xdr  float64 `json:"xdr"`
			Xlm  float64 `json:"xlm"`
			Xrp  float64 `json:"xrp"`
			Yfi  float64 `json:"yfi"`
			Zar  float64 `json:"zar"`
			Bits float64 `json:"bits"`
			Link float64 `json:"link"`
			Sats float64 `json:"sats"`
		} `json:"current_price"`
		MarketCap struct {
			Aed  float64 `json:"aed"`
			Ars  float64 `json:"ars"`
			Aud  float64 `json:"aud"`
			Bch  float64 `json:"bch"`
			Bdt  float64 `json:"bdt"`
			Bhd  float64 `json:"bhd"`
			Bmd  float64 `json:"bmd"`
			Bnb  float64 `json:"bnb"`
			Brl  float64 `json:"brl"`
			Btc  float64 `json:"btc"`
			Cad  float64 `json:"cad"`
			Chf  float64 `json:"chf"`
			Clp  float64 `json:"clp"`
			Cny  float64 `json:"cny"`
			Czk  float64 `json:"czk"`
			Dkk  float64 `json:"dkk"`
			Dot  float64 `json:"dot"`
			Eos  float64 `json:"eos"`
			Eth  float64 `json:"eth"`
			Eur  float64 `json:"eur"`
			Gbp  float64 `json:"gbp"`
			Hkd  float64 `json:"hkd"`
			Huf  float64 `json:"huf"`
			Idr  float64 `json:"idr"`
			Ils  float64 `json:"ils"`
			Inr  float64 `json:"inr"`
			Jpy  float64 `json:"jpy"`
			Krw  float64 `json:"krw"`
			Kwd  float64 `json:"kwd"`
			Lkr  float64 `json:"lkr"`
			Ltc  float64 `json:"ltc"`
			Mmk  float64 `json:"mmk"`
			Mxn  float64 `json:"mxn"`
			Myr  float64 `json:"myr"`
			Ngn  float64 `json:"ngn"`
			Nok  float64 `json:"nok"`
			Nzd  float64 `json:"nzd"`
			Php  float64 `json:"php"`
			Pkr  float64 `json:"pkr"`
			Pln  float64 `json:"pln"`
			Rub  float64 `json:"rub"`
			Sar  float64 `json:"sar"`
			Sek  float64 `json:"sek"`
			Sgd  float64 `json:"sgd"`
			Thb  float64 `json:"thb"`
			Try  float64 `json:"try"`
			Twd  float64 `json:"twd"`
			Uah  float64 `json:"uah"`
			Usd  float64 `json:"usd"`
			Vef  float64 `json:"vef"`
			Vnd  float64 `json:"vnd"`
			Xag  float64 `json:"xag"`
			Xau  float64 `json:"xau"`
			Xdr  float64 `json:"xdr"`
			Xlm  float64 `json:"xlm"`
			Xrp  float64 `json:"xrp"`
			Yfi  float64 `json:"yfi"`
			Zar  float64 `json:"zar"`
			Bits float64 `json:"bits"`
			Link float64 `json:"link"`
			Sats float64 `json:"sats"`
		} `json:"market_cap"`
		TotalVolume struct {
			Aed  float64 `json:"aed"`
			Ars  float64 `json:"ars"`
			Aud  float64 `json:"aud"`
			Bch  float64 `json:"bch"`
			Bdt  float64 `json:"bdt"`
			Bhd  float64 `json:"bhd"`
			Bmd  float64 `json:"bmd"`
			Bnb  float64 `json:"bnb"`
			Brl  float64 `json:"brl"`
			Btc  float64 `json:"btc"`
			Cad  float64 `json:"cad"`
			Chf  float64 `json:"chf"`
			Clp  float64 `json:"clp"`
			Cny  float64 `json:"cny"`
			Czk  float64 `json:"czk"`
			Dkk  float64 `json:"dkk"`
			Dot  float64 `json:"dot"`
			Eos  float64 `json:"eos"`
			Eth  float64 `json:"eth"`
			Eur  float64 `json:"eur"`
			Gbp  float64 `json:"gbp"`
			Hkd  float64 `json:"hkd"`
			Huf  float64 `json:"huf"`
			Idr  float64 `json:"idr"`
			Ils  float64 `json:"ils"`
			Inr  float64 `json:"inr"`
			Jpy  float64 `json:"jpy"`
			Krw  float64 `json:"krw"`
			Kwd  float64 `json:"kwd"`
			Lkr  float64 `json:"lkr"`
			Ltc  float64 `json:"ltc"`
			Mmk  float64 `json:"mmk"`
			Mxn  float64 `json:"mxn"`
			Myr  float64 `json:"myr"`
			Ngn  float64 `json:"ngn"`
			Nok  float64 `json:"nok"`
			Nzd  float64 `json:"nzd"`
			Php  float64 `json:"php"`
			Pkr  float64 `json:"pkr"`
			Pln  float64 `json:"pln"`
			Rub  float64 `json:"rub"`
			Sar  float64 `json:"sar"`
			Sek  float64 `json:"sek"`
			Sgd  float64 `json:"sgd"`
			Thb  float64 `json:"thb"`
			Try  float64 `json:"try"`
			Twd  float64 `json:"twd"`
			Uah  float64 `json:"uah"`
			Usd  float64 `json:"usd"`
			Vef  float64 `json:"vef"`
			Vnd  float64 `json:"vnd"`
			Xag  float64 `json:"xag"`
			Xau  float64 `json:"xau"`
			Xdr  float64 `json:"xdr"`
			Xlm  float64 `json:"xlm"`
			Xrp  float64 `json:"xrp"`
			Yfi  float64 `json:"yfi"`
			Zar  float64 `json:"zar"`
			Bits float64 `json:"bits"`
			Link float64 `json:"link"`
			Sats float64 `json:"sats"`
		} `json:"total_volume"`
	} `json:"market_data"`
	CommunityData struct {
		FacebookLikes            interface{} `json:"facebook_likes"`
		TwitterFollowers         int         `json:"twitter_followers"`
		RedditAveragePosts48H    float64     `json:"reddit_average_posts_48h"`
		RedditAverageComments48H float64     `json:"reddit_average_comments_48h"`
		RedditSubscribers        int         `json:"reddit_subscribers"`
		RedditAccountsActive48H  string      `json:"reddit_accounts_active_48h"`
	} `json:"community_data"`
	DeveloperData struct {
		Forks                        interface{} `json:"forks"`
		Stars                        interface{} `json:"stars"`
		Subscribers                  interface{} `json:"subscribers"`
		TotalIssues                  interface{} `json:"total_issues"`
		ClosedIssues                 interface{} `json:"closed_issues"`
		PullRequestsMerged           interface{} `json:"pull_requests_merged"`
		PullRequestContributors      interface{} `json:"pull_request_contributors"`
		CodeAdditionsDeletions4Weeks struct {
			Additions interface{} `json:"additions"`
			Deletions interface{} `json:"deletions"`
		} `json:"code_additions_deletions_4_weeks"`
		CommitCount4Weeks interface{} `json:"commit_count_4_weeks"`
	} `json:"developer_data"`
	PublicInterestStats struct {
		AlexaRank   int         `json:"alexa_rank"`
		BingMatches interface{} `json:"bing_matches"`
	} `json:"public_interest_stats"`
}

type Config struct {
	NotaryID []string `json:"notary_id"`
	Range    []string `json:"range"`
	Api      []string `json:"api"`
	Fiat     []string `json:"fiat"`
	Years    []int    `json:"years"`
}

type SmkStatsKmdMining struct {
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

type KomodoStatsKmdMining []struct {
	Name     string `json:"name"`
	Iguanaid int    `json:"iguanaid"`
	Season   int    `json:"season"`
	Region   string `json:"region"`
	Height   int    `json:"height"`
	Address  string `json:"address"`
	Amount   string `json:"amount"`
	Time     int    `json:"time"`
	Txid     string `json:"txid"`
	Date     string
	FiatUSD  float64
	FiatEUR  float64
	TotalKMD string
}

var gConfig Config
var gSmkStatsMining SmkStatsKmdMining
var gKomodoStatsMining KomodoStatsKmdMining
var gGeckoRegistry = make(map[string]CoingeckoHistoryResponse)
var gCurYear string

func InternalExecGet(finalEndpoint string, ctx *fasthttp.RequestCtx, shouldRelease bool) (*fasthttp.Request, *fasthttp.Response) {
	_ = glg.Debugf("final endpoint: %s", finalEndpoint)
	status, body, err := fasthttp.GetTimeout(nil, finalEndpoint, 30*time.Second)
	if err != nil {
		_ = glg.Errorf("Http error: %v & endpoint: %s", err, finalEndpoint)
	}
	if ctx != nil {
		ctx.SetStatusCode(status)
		ctx.SetBodyString(string(body))
	}
	if len(string(body)) > 100 {
		_ = glg.Debugf("http response (100 first bytes): %s", string(body)[0:100])
	} else {
		_ = glg.Debugf("http response: %s", string(body))
	}
	if !shouldRelease {
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		res.SetStatusCode(status)
		res.SetBody(body)
		return req, res
	}
	return nil, nil
}

func ReleaseInternalExecGet(req *fasthttp.Request, res *fasthttp.Response) {
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(res)
}

func parseConfig() {
	file, _ := ioutil.ReadFile("config.json")
	_ = json.Unmarshal([]byte(file), &gConfig)
	fmt.Printf("%v\n", gConfig)
}

func main() {
	parseConfig()
	for _, s := range gConfig.NotaryID {
		for _, curYear := range gConfig.Years {
			fmt.Printf("Calculating taxe revenue of %s\n", s)
			for _, curRange := range gConfig.Range {
				switch curRange {
				case "month":
					calculateTaxeLastMonth(s)
					break
				case "year":
					calculateTaxeLastYear(s)
					break
				case "1":
					calculateBySpecificMonth(s, time.January, curYear)
					break
				case "2":
					calculateBySpecificMonth(s, time.February, curYear)
					break
				case "3":
					calculateBySpecificMonth(s, time.March, curYear)
					break
				case "4":
					calculateBySpecificMonth(s, time.April, curYear)
					break
				case "5":
					calculateBySpecificMonth(s, time.May, curYear)
					break
				case "6":
					calculateBySpecificMonth(s, time.June, curYear)
					break
				case "7":
					calculateBySpecificMonth(s, time.July, curYear)
					break
				case "8":
					calculateBySpecificMonth(s, time.August, curYear)
					break
				case "9":
					calculateBySpecificMonth(s, time.September, curYear)
					break
				case "10":
					calculateBySpecificMonth(s, time.October, curYear)
					break
				case "11":
					calculateBySpecificMonth(s, time.November, curYear)
					break
				case "12":
					calculateBySpecificMonth(s, time.December, curYear)
					break
				}
			}
		}
	}
}

func fillMiningStats(notaryNodeId string, from time.Time, to time.Time, category string) {
	for _, curApi := range gConfig.Api {
		switch curApi {
		case "https://stats.kmd.io":
			fillSmkStatsMining(notaryNodeId, strconv.FormatInt(from.Unix(), 10), strconv.FormatInt(to.Unix(), 10), category)
			break
		case "https://komodostats.com/":
			fillKomodoStatsMining(notaryNodeId, from.Format(layoutISO), to.Format(layoutISO), category)
			break
		}
	}
}

func fillKomodoStatsMining(notaryNodeId string, from string, to string, category string) {
	url := komodoStatsUri + notaryNodeId + "&start=" + from + "&end=" + to
	fmt.Printf("Processing: %s\n", url)
	req, res := InternalExecGet(url, nil, false)
	_ = glg.Infof("%d", res.StatusCode())
	if res.StatusCode() == 200 {
		err := json.Unmarshal(res.Body(), &gKomodoStatsMining)
		if err != nil {
			return
		}
	}
	ReleaseInternalExecGet(req, res)
	CalculateFromKomodoStatsMining(category, notaryNodeId)
}

func CalculateFromKomodoStatsMining(category string, notaryNodeId string) {
	var prevValueUSD float64 = 0
	var prevValueEUR float64 = 0
	var totalKMD float64 = 0
	for i := range gKomodoStatsMining {
		curTransaction := gKomodoStatsMining[len(gKomodoStatsMining)-1-i]
		//fmt.Printf("%d\n", curTransaction.Time)
		var curTime = time.Unix(int64(curTransaction.Time), 0)
		curTransaction.Date = curTime.Format(layoutOutCsv)
		fmt.Printf("%d - %s - %s\n", curTransaction.Height, curTime.Format(layoutCoingecko), curTime.Format(layoutOutCsv))
		from := curTime.Format(layoutCoingecko)
		processCoingecko(from)
		if amountF, err := strconv.ParseFloat(curTransaction.Amount, 64); err == nil {
			if i > 0 {
				//fmt.Printf("%v\n", gKomodoStatsMining[len(gKomodoStatsMining)-i])
				prevValueUSD = gKomodoStatsMining[len(gKomodoStatsMining)-i].FiatUSD
				prevValueEUR = gKomodoStatsMining[len(gKomodoStatsMining)-i].FiatEUR
			}
			totalKMD += amountF
			curTransaction.TotalKMD = fmt.Sprintf("%f", totalKMD)
			curTransaction.FiatUSD = prevValueUSD + (amountF * gGeckoRegistry[from].MarketData.CurrentPrice.Usd)
			curTransaction.FiatEUR = prevValueEUR + (amountF * gGeckoRegistry[from].MarketData.CurrentPrice.Eur)
			gKomodoStatsMining[len(gKomodoStatsMining)-1-i] = curTransaction
			//fmt.Printf("value in usd: cur: %.2f prev: %.2f\n", curTransaction.FiatUSD, prevValueUSD)
			//fmt.Printf("value in eur: %.2f prev: %.2f\n", curTransaction.FiatEUR, prevValueEUR)
		} else {
			fmt.Printf("What the fuck")
		}
	}
	generateKomodoStatsTransactionReport(category, notaryNodeId)
}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func generateKomodoStatsTransactionReport(category string, notaryNodeId string) {
	file, err := os.Create(notaryNodeId + "-" + category + "-komodostats-mining-tax-report-" + gCurYear + ".csv")
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var headers = []string{"Name", "Address", "TransactionUrl", "BlockUrl", "Date", "Amount", "TotalKMD", "TotalUSD", "TotalEUR"}
	err = writer.Write(headers)
	checkError("Cannot write to file", err)
	for i := range gKomodoStatsMining {
		curTransaction := gKomodoStatsMining[len(gKomodoStatsMining)-1-i]
		usdAmount := fmt.Sprintf("%.2f", curTransaction.FiatUSD)
		eurAmount := fmt.Sprintf("%.2f", curTransaction.FiatEUR)
		var result = []string{curTransaction.Name, curTransaction.Address, explorerTxUri + curTransaction.Txid, explorerBlockUri + strconv.Itoa(curTransaction.Height), curTransaction.Date, curTransaction.Amount, curTransaction.TotalKMD, usdAmount, eurAmount}
		err = writer.Write(result)
		checkError("Cannot write to file", err)
	}
}

func processCoingecko(from string) {
	_, ok := gGeckoRegistry[from]
	if ok {
	} else {
		url := coingeckoUri + from + "&localization=false"
		fmt.Printf("%s not found, processing: %s\n", from, url)
		req, res := InternalExecGet(url, nil, false)
		_ = glg.Infof("%d", res.StatusCode())
		var geckoAnswer = CoingeckoHistoryResponse{}
		if res.StatusCode() == 200 {
			err := json.Unmarshal(res.Body(), &geckoAnswer)
			if err != nil {
				return
			}
			gGeckoRegistry[from] = geckoAnswer
		} else if res.StatusCode() == 429 {
			fmt.Printf("429 limit exceeded, waiting one second before continuing\n")
			time.Sleep(1 * time.Second)
			processCoingecko(from)
		}
		ReleaseInternalExecGet(req, res)
	}

}

func fillSmkStatsMining(notaryNodeId string, from string, to string, category string) {
	url := smkStatsUri + notaryNodeId + "&min_blocktime=" + from + "&max_blocktime=" + to
	fmt.Printf("Processing: %s\n", url)

}

func calculateTaxeLastYear(notaryNodeId string) {
	fmt.Printf("Calculating tax for last year: [%s]\n", notaryNodeId)
	beginningOfYear := now.BeginningOfYear()
	endOfYear := now.EndOfYear()
	fmt.Printf("timestamp begin of the year: %v %d\n", beginningOfYear, beginningOfYear.Unix())
	fmt.Printf("timestamp end of the year: %v %d\n", endOfYear, endOfYear.Unix())
	fillMiningStats(notaryNodeId, beginningOfYear, endOfYear, "yearly")
}

func calculateTaxeLastMonth(notaryNodeId string) {
	fmt.Printf("Calculating tax for last month: [%s]\n", notaryNodeId)
	beginningOfMonth := now.BeginningOfMonth()
	endOfMonth := now.EndOfMonth()
	fmt.Printf("timestamp begin of the month: %v %d\n", beginningOfMonth, beginningOfMonth.Unix())
	fmt.Printf("timestamp end of the month: %v %d\n", endOfMonth, endOfMonth.Unix())
	fillMiningStats(notaryNodeId, beginningOfMonth, endOfMonth, "monthly")
}

func calculateBySpecificMonth(notaryNodeId string, month time.Month, year int) {
	var deducedTime time.Time = time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	beginningOfMonth := now.With(deducedTime).BeginningOfMonth()
	endOfMonth := now.With(deducedTime).EndOfMonth()
	fmt.Printf("timestamp begin of the month: %v %d\n", beginningOfMonth, beginningOfMonth.Unix())
	fmt.Printf("timestamp end of the month: %v %d\n", endOfMonth, endOfMonth.Unix())
	gCurYear = strconv.Itoa(year)
	fillMiningStats(notaryNodeId, beginningOfMonth, endOfMonth, month.String())
}
