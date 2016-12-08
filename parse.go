package koreannet

import (
	"database/sql"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)


const (
	DB_HOST = "tcp(127.0.0.1:3306)"
	DB_NAME = "database"
	DB_USER = "user"
	DB_PASS = "password"
)

type Info struct {
	XMLName             xml.Name  `xml:"MGS1OutXml"`
	Gtin                string    `xml:"gtin"`
	Dscrgtink           string    `xml:"dscrgtink"`
	Npname              string    `xml:"npname"`
	Conamek             string    `xml:"conamek"`
	Conamee             string    `xml:"conamee"`
	Dscrbrandk          string    `xml:"dscrbrandk"`
	Countrydescr        string    `xml:"countrydescr"`
	Dstartavailble      string    `xml:"dstartavailble"`
	Dsysupdated         string    `xml:"dsysupdated"`
	Imgpath1            string    `xml:"imgpath1"`
	Imgpath2            string    `xml:"imgpath2"`
	Imgpath3            string    `xml:"imgpath3"`
	Pgurl               string    `xml:"pgurl"`
	Detail_text         string    `xml:"text"`
	Kanclasscode        string    `xml:"kanclasscode"`
	Unitnetcont         string    `xml:"unitnetcont"`
	Unitnetcontuomdescr string    `xml:"unitnetcontuomdescr"`
	Unitnetcontuom      string    `xml:"unitnetcontuom"`
	Unitsinpack         string    `xml:"unitsinpack"`
	Height              string    `xml:"height"`
	Heightuomdescr      string    `xml:"heightuomdescr"`
	Heightuom           string    `xml:"heightuom"`
	Width               string    `xml:"width"`
	Widthuomdescr       string    `xml:"widthuomdescr"`
	Widthuom            string    `xml:"widthuom"`
	Depth               string    `xml:"depth"`
	Depthuomdescr       string    `xml:"depthuomdescr"`
	Depthuom            string    `xml:"depthuom"`
	Netweight           string    `xml:"netweight"`
	Netweightuomdescr   string    `xml:"netweightuomdescr"`
	Netweightuom        string    `xml:"netweightuom"`
	Grossweight         string    `xml:"grossweight"`
	Grossweightuomdescr string    `xml:"grossweightuomdescr"`
	Grossweightuom      string    `xml:"grossweightuom"`
	Busnid              string    `xml:"busnid"`
	Producturl          string    `xml:"producturl"`
	Datasource          string    `xml:"datasource"`
	Prgubun             string    `xml:"prgubun"`
	Packgtin            string    `xml:"packgtin"`
	Imgpath4            string    `xml:"imgpath4"`
	Product             []Product `xml:"PRODUCT"`
}

type Product struct {
	XMLName xml.Name `xml:"PRODUCT"`
	Login   string   `xml:"login"`
}

func Parse(wg *sync.WaitGroup, seq int, barcode string, id string, pw string) bool {
	defer wg.Done()

	var result_flag bool
	dsn := DB_USER + ":" + DB_PASS + "@" + DB_HOST + "/" + DB_NAME
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    
        if err = db.Ping(); err != nil {
		defer func() {
			fmt.Println(err)
			// if err := recover(); err != nil {
			// 	log.Fatal(err)
			// }
		}()
		return false
	}

	fmt.Printf("Process Start Id : %v (%v)\n", seq, time.Now())

	path := fmt.Sprintf("http://api.koreannet.or.kr/mobileweb/search/barcodeSearchXml.do?barcode=%s&id=%s&pw=%s", barcode, id, pw)
	response, err := http.Get(path)

	/***** if post method *****/
	path = "http://api.koreannet.or.kr/mobileweb/search/barcodeSearchXml.do"
	urlData := url.Values{}
	urlData.Set("barcode", barcode)
	urlData.Set("id", id)
	urlData.Set("pw", pw)
	response, err := http.PostForm(path, urlData)
	/**************************/

	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	content := string(data)

	q := Info{}
	xml.Unmarshal([]byte(content), &q)

	var result string
	for _, res := range q.Product {
		result = res.Login
	}

	re, _ := regexp.Compile(`오류`)
	res := re.FindAllStringSubmatch(result, -1)

	if len(res) > 0 {
		fmt.Println(result)
		result_flag = false
	} else {

	    defer func() {
			if err := recover(); err != nil {
			    logFile, _ := os.OpenFile("/tmp/koreannet.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
			    defer logFile.Close()

			    log.SetOutput(logFile)
			    log.Println(err)
			}
		}()

		sql := "insert into barcode_product_info set ID=?, PW=?, GTIN=?, DSCRGTINK=?, NPNAME=?"
		sql = fmt.Sprintf("%s,CONAMEK=?,CONAMEE=?,DSCRBRANDK=?,COUNTRYDESCR=?", sql)
		sql = fmt.Sprintf("%s,DSTARTAVAILBLE=?,DSYSUPDATED=?,IMGPATH1=?,IMGPATH2=?", sql)
		sql = fmt.Sprintf("%s,IMGPATH3=?,IMGPATH4=?,PGURL=?,DETAIL_TEXT=?", sql)
		sql = fmt.Sprintf("%s,KANCLASSCODE=?,UNITNETCONT=?,UNITNETCONTUOMDESCR=?,UNITNETCONTUOM=?", sql)
		sql = fmt.Sprintf("%s,UNITSINPACK=?,HEIGHT=?,HEIGHTUOMDESCR=?,HEIGHTUOM=?,WIDTH=?", sql)
		sql = fmt.Sprintf("%s,WIDTHUOMDESCR=?,WIDTHUOM=?,DEPTH=?,DEPTHUOMDESCR=?", sql)
		sql = fmt.Sprintf("%s,DEPTHUOM=?,NETWEIGHT=?,NETWEIGHTUOMDESCR=?,NETWEIGHTUOM=?", sql)
		sql = fmt.Sprintf("%s,GROSSWEIGHT=?,GROSSWEIGHTUOMDESCR=?,GROSSWEIGHTUOM=?", sql)
		sql = fmt.Sprintf("%s,BUSNID=?,PRODUCTURL=?,DATASOURCE=?,PRGUBUN=?,PACKGTIN=?", sql)

		stmt, _ := db.Prepare(sql)
		defer stmt.Close()
		_, err := stmt.Exec(id, pw, q.Gtin, q.Dscrgtink, q.Npname, q.Conamek, q.Conamee, q.Dscrbrandk, q.Countrydescr, q.Dstartavailble, q.Dsysupdated, q.Imgpath1, q.Imgpath2, q.Imgpath3, q.Imgpath4, q.Pgurl, q.Detail_text, q.Kanclasscode, q.Unitnetcont, q.Unitnetcontuomdescr, q.Unitnetcontuom, q.Unitsinpack, q.Height, q.Heightuomdescr, q.Heightuom, q.Width, q.Widthuomdescr, q.Widthuom, q.Depth, q.Depthuomdescr, q.Depthuom, q.Netweight, q.Netweightuomdescr, q.Netweightuom, q.Grossweight, q.Grossweightuomdescr, q.Grossweightuom, q.Busnid, q.Producturl, q.Datasource, q.Prgubun, q.Packgtin)

		if err != nil {
			log.Fatal(err)
		}

		result_flag = true
	}

	fmt.Printf("Process Stopp Id : %v (%v)\n", seq, time.Now())

	return result_flag
}
