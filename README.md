# 코리안넷 상품정보 조회



### 단위작업
```go
package main

import "github.com/jhjang-dev/koreannet"
import "sync"
const (
    code = "1234567890"
    id = "test"
    pw = "pass"
)

var wg sync.WaitGroup

func main(){
    seq := 1
    wg.Add(1)
    result := Parse(&wg,seq,code,id,pw)
    wg.Wait()
}
```

### 데몬방식
```go
package main
import (
	"sync"

    "github.com/jhjang-dev/koreannet"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

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
			"message": "pong",
			"barcode": code,
			"id":      id,
			"pw":      pw,
		})
	})
	wg.Wait()
	r.Run(":9005") // listen and server on 0.0.0.0:8080
}
```

### Example
```go
http://localhost:9005?barcode=91827364505729&id=test&pw=1234
```
