package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"sync"
	"time"
)

type NotifyRobot struct {
	mLastSendTimeTable sync.Map
}

func (this *NotifyRobot) getClientAddr(ctx *gin.Context) string {

	dstData := ctx.ClientIP()

	var srcIp string

	//Cf-Connecting-Ip
	srcIp = ctx.GetHeader("CF-Connecting-IP")
	if len(srcIp) > 0 {
		dstData = srcIp
		return dstData
	}

	//X-Real-Ip
	srcIp = ctx.GetHeader("X-Real-IP")
	if len(srcIp) > 0 {
		dstData = srcIp
		return dstData
	}

	//X-Forwarded-For
	srcIp = ctx.GetHeader("X-Forwarded-For")
	if len(srcIp) > 0 {
		dstData = srcIp
		return dstData
	}

	return dstData

}

func (this *NotifyRobot) getWebHookUrlForFeishu() string {

	var dstUrl string

	dstUrl = os.Getenv("NOTIFY_WEBHOOK_FEISHU")

	return dstUrl

}

func (this *NotifyRobot) sendNotifyFeishu(ctx *gin.Context) error {

	var xErr error

	xWebHookUrl := this.getWebHookUrlForFeishu()
	if len(xWebHookUrl) < 1 {
		return xErr
	}

	log.Infof("设置了飞书Webhook回调=[%v]", xWebHookUrl)

	xNowTime := time.Now().Format("2006-01-02 15:04:05")

	xClientAddr := this.getClientAddr(ctx)

	xJsonData := `{
		"msg_type": "text",
		"content": {
			"text": "【通知】[` + xNowTime + `][` + xClientAddr + `]触发了[OpenList]访问"
		}
	}`

	xHttpReq, xHttpErr := http.NewRequest("POST", xWebHookUrl, bytes.NewBuffer([]byte(xJsonData)))

	if xHttpErr != nil {
		return xErr
	}

	xHttpReq.Header.Set("Content-Type", "application/json")
	xHttpReq.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36")

	xHttpClient := &http.Client{}

	xHttpResp, xHttpErr := xHttpClient.Do(xHttpReq)
	if xHttpErr != nil {
		return xErr
	}

	defer xHttpResp.Body.Close()

	return xErr

}

func (this *NotifyRobot) Notify(ctx *gin.Context) error {

	var xErr error

	xClientAddr := ctx.ClientIP()
	var xClientLastTime int64 = 0

	xLastTimeVal, xIsLoad := this.mLastSendTimeTable.Load(xClientAddr)
	if xIsLoad {
		xClientLastTime = xLastTimeVal.(int64)
	}

	xTimeDiff := time.Now().UnixMilli() - xClientLastTime

	log.Infof("NotifyRobot xClientAddr=[%v] xClientLastTime=[%v] xTimeDiff=[%v]", xClientAddr, xClientLastTime, xTimeDiff/1000)

	this.mLastSendTimeTable.Store(xClientAddr, time.Now().UnixMilli())

	if xTimeDiff < 5*60*1000 {
		return xErr
	}

	this.sendNotifyFeishu(ctx)

	return xErr

}

var (
	gNotifyRobot NotifyRobot
)

func AccessNotify(c *gin.Context) {

	gNotifyRobot.Notify(c)

	c.Next()

}
