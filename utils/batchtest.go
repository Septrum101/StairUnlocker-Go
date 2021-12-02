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
	for i := range proxiesList {
		p := proxiesList[i]
		b.Go(p.Name(), func() (interface{}, error) {
			latency, sCode, err := NETFLIXTest(p, "https://www.netflix.com/title/70143836")
			if err != nil {
				log.Errorln("%s : %s", p.Name(), err.Error())
			} else if sCode == 200 {
				log.Infoln("%s : latency = %v ms | Full Unlock", p.Name(), latency)
				NETFLIXList = append(NETFLIXList, p.Name())
			} else {
				log.Infoln("%s : latency = %v ms | None", p.Name(), latency)
			}
			return nil, nil
		})
	}
	b.Wait()
	return
}
