package utils

import (
	"context"
	"github.com/Dreamacro/clash/common/batch"
	C "github.com/Dreamacro/clash/constant"
	"github.com/Dreamacro/clash/log"
)

//BatchCheck : n int, to set ConcurrencyNum.
func BatchCheck(proxiesList []C.Proxy, n int) (NETFLIXList []string) {
	b, _ := batch.New(context.Background(), batch.WithConcurrencyNum(n))
	curr, total := 0, len(proxiesList)
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			latency, sCode, err := NETFLIXTest(p, "https://www.netflix.com/title/70143836")
			if err != nil {
				curr++
				log.Errorln("(%d/%d) %s : %s", curr, total, p.Name(), err.Error())
			} else if sCode == 200 {
				curr++
				log.Infoln("(%d/%d) %s : latency = %v ms | Full Unlock", curr, total, p.Name(), latency)
				NETFLIXList = append(NETFLIXList, p.Name())
			} else {
				curr++
				log.Infoln("(%d/%d) %s : latency = %v ms | None", curr, total, p.Name(), latency)
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}
