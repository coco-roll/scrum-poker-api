package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"math/rand"
	"net/http"
	"scrum-poker/models"
	"scrum-poker/pkg/e"
	"strconv"
	"strings"
)

//查
func GetTest(c *gin.Context) {
	str := c.Query("id")
	id, err := strconv.Atoi(str)
	if err != nil {
		code := e.INVALID_PARAMS
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": "",
		})
		models.CloseDB()
		return
	}
	data := make(map[string]interface{})

	code := e.SUCCESS

	data["model"] = models.GetTestModel(id)

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//增
func AddTest(c *gin.Context) {
}

//改
func EditTest(c *gin.Context) {
}

//删
func DeleteTest(c *gin.Context) {
}

//得到链接
func GetUrl(c *gin.Context) {
	randNumber := rand.Int()
	data := make(map[string]interface{})
	data["code"] = strconv.Itoa(randNumber)
	returnjson(c, 1, data)
}

//翻牌
func Poker(c *gin.Context) {
	var spoker models.Scrum_poker

	spoker.Url_code = c.Query("url_code")
	spoker.Poker = c.Query("poker")
	user_id, err := c.Cookie("user_id")
	if err != nil {
		user_id = "1"
	}
	spoker.User_id, _ = strconv.Atoi(user_id)

	where := make(map[string]interface{})
	where["user_id"] = spoker.User_id
	where["Url_code"] = spoker.Url_code

	info := models.GetOne(where)
	if info.Id == 0 {
		models.AddPoker(spoker)
	} else {
		models.UpdPoker(where, spoker)
	}

	data := make(map[string]interface{})
	data["msg"] = "成功"
	returnjson(c, 1, data)
}

//设置cookier
func SetCk(c *gin.Context) {
	_, err := c.Cookie("user_id")
	if err != nil {
		randNumber := rand.Intn(10000)
		c.SetCookie("user_id", strconv.Itoa(randNumber), 3600, "/", "", false, true)
	}
}

//返回值
func returnjson(c *gin.Context, status int, data map[string]interface{}) {

	code := e.SUCCESS
	if status == 0 {
		code = e.ERROR
	} else if status == 2 {
		code = e.INVALID_PARAMS
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}

//websocket
var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsPoker struct {
	Ws    *websocket.Conn
	Poker string `json:"poker"`
	Code  string
	Mt    int
}

var Wspokers map[string][]WsPoker

//webSocket请求ping 返回pong
func Ping(c *gin.Context) {
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()
	code := ""
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		msg := string(message)

		//链接
		if strings.Index(msg, "type=1") != -1 {
			if len(Wspokers) == 0 {
				Wspokers = make(map[string][]WsPoker)
			}
			code = getUrlParams(msg, "code")
			fmt.Println(code)
			if code == "" {
				fmt.Println("群组不存在[" + msg + "]")
				break
			}
			wsclient := WsPoker{Ws: ws, Code: code, Mt: mt}
			Wspokers[code] = append(Wspokers[code], wsclient)
			fmt.Println(Wspokers)
			//翻牌
		} else if strings.Index(msg, "type=2") != -1 {
			poker := getUrlParams(msg, "poker")
			if poker == "" {
				fmt.Println("选项卡不存在[" + msg + "]")
				break
			}
			//var wsPoker WsPoker
			for k, v := range Wspokers[code] {
				if v.Ws == ws {
					Wspokers[code][k].Poker = poker
				}
			}
			fmt.Println(Wspokers)
			//重新开始
		} else if strings.Index(msg, "type=3") != -1 {
			for k, _ := range Wspokers[code] {
				Wspokers[code][k].Poker = ""
			}
			fmt.Println(Wspokers[code])
		}
	}
	//断开链接
	for k, v := range Wspokers[code] {
		if v.Ws == ws {
			Wspokers[code] = append(Wspokers[code][:k], Wspokers[code][k+1:]...)
		}
	}
	fmt.Println(Wspokers[code])
}

func getUrlParams(url, key string) string {
	strArr := strings.Split(url, "&")
	for _, v := range strArr {
		valArr := strings.Split(v, "=")

		if key == valArr[0] {
			return valArr[1]
		}

	}
	return ""
}
