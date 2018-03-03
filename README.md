# sheep


## example

```
import (
	"log"

	"time"

	"github.com/leek-box/sheep/huobi"
)

func main() {
	h, err := huobi.NewHuobi("your-access-key", "your-secret-key")
	if err != nil {
		log.Println(err.Error())
		return
	}

	h.GetAccountBalance()

	listener := func(symbol string, depth *huobi.MarketDepth) {
		log.Println(depth)
	}
	h.SetDepthlListener(listener)

	h.SubscribeDepth("btcusdt")
	h.SubscribeDepth("ethusdt")

	time.Sleep(time.Hour)
}

```
