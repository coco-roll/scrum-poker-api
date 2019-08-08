package main

import (
	"fmt"
	"time"
	"strconv"
	"net/http"
	"encoding/json"
	"scrum-poker/pkg/setting"
	"scrum-poker/routers"
	"scrum-poker/app/api"
)

func isOverNot(key string) (is_over bool) {
    // 玩家个数
    sec, _:= setting.Cfg.GetSection("app")
    palyer_num, _:= strconv.Atoi(sec.Key("PLAYER_NUM").String())
    fmt.Println(palyer_num)
    len_wspokers := len(api.Wspokers[key])
    if (len_wspokers == 0) {
        return false
    }

    for _, v := range api.Wspokers[key] {
        if ((len(v.Poker) == 0) || len_wspokers == 1) {
            return false
        }
    }

    return true
}

func needleTrack() {
    ticker:=time.NewTicker(time.Second*2)
    go func() {
        for _=range ticker.C {
            for  k, _:= range api.Wspokers {
                // 游戏是否结束
                if (isOverNot(k)) {
                    res := make(map[string]int)
                    for  _,v := range api.Wspokers[k]{
                        
                        if _, ok := res[v.Poker]; ok {
                            res[v.Poker] += 1
                        }else{
                            res[v.Poker] = 1
                        }
                        
                    }

                    mjson,_ :=json.Marshal(res)
                    // 写入ws数据
                    for  _, v := range api.Wspokers[k]{
                        err := v.Ws.WriteMessage(v.Mt, mjson)
                        if err != nil {
                            fmt.Println("[" + strconv.Itoa(v.Mt) + "]:推送结果失败")
                        }
                    }
                    //断开链接
                    // delete(api.Wspokers, k)
                }
            }
        }
    }()
   time.Sleep(time.Minute)
}

func main() {
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// 探针监测游戏是否结束
    go needleTrack()

	s.ListenAndServe()
}
