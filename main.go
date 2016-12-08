
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var wg sync.WaitGroup

func main() {
	mode := flag.String("mode", "single or daemon", "실행형식")

	flag.Parse()

	if flag.NFlag() == 0 {
		flag.Usage()
		return
	}

	if *mode == "daemon" {
		r := gin.Default()
		r.GET("/search", func(c *gin.Context) {

			code := c.Query("barcode")
			id := c.Query("id")
			pw := c.Query("pw")

			seq := 1

			var result bool
			wg.Add(1)
			go func(wg sync.WaitGroup, seq int, code string, id string, pw string) {
				result = Parse(&wg, seq, code, id, pw)
			}(wg, seq, code, id, pw)

			c.JSON(200, gin.H{
				"barcode": code,
				"id":      id,
				"pw":      pw,
				"result":  result,
			})
		})
		wg.Wait()
		r.Run(":9005") // listen and server on 0.0.0.0:8080
	} else {

		runtime.GOMAXPROCS(runtime.NumCPU())

		if len(os.Args) < 3 {
			fmt.Println("Usage: file -mode=single barcode^^id^^password")
			return
		}

		barcodes := os.Args[2:]

		seq := 1
		for _, dt := range barcodes {

			result := strings.Split(dt, "^^")
			arr_len := len(result)

			if arr_len < 3 {

				fmt.Println("Format is incorrect")
				continue
			}

			code := result[0]
			id := result[1]
			pw := result[2]

			wg.Add(1)
			go Parse(&wg, seq, code, id, pw)

			if seq%10 == 0 {
				fmt.Println("wait:", seq)
				wg.Wait()
			}

			seq += 1
		}
		wg.Wait()

	}

}
