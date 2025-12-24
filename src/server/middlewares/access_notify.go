package middlewares

import (
	"fmt"
	"github.com/OpenListTeam/OpenList/v4/pkg/webhook"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
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

func (this *NotifyRobot) Notify(ctx *gin.Context) error {

	var dstErr error

	if !webhook.WebHookIsSet() {
		return dstErr
	}

	clientIP := ctx.ClientIP()

	var clientLastTime int64 = 0

	lastTimeValue, lastTimeOk := this.mLastSendTimeTable.Load(clientIP)
	if lastTimeOk {
		clientLastTime = lastTimeValue.(int64)
	}

	lastTimeDiff := time.Now().UnixMilli() - clientLastTime

	log.Infof("[NotifyRobot]clientIP=[%v] clientLastTime=[%v] lastTimeDiff=[%v]", clientIP, clientLastTime, lastTimeDiff/1000)

	this.mLastSendTimeTable.Store(clientIP, time.Now().UnixMilli())

	if lastTimeDiff < 5*60*1000 {
		return dstErr
	}

	msgTime := time.Now().Format("2006-01-02 15:04:05")

	msgData2Feishu := fmt.Sprintf("【通知】[%v]IP地址=[%v]触发了[OpenList]访问", msgTime, clientIP)

	dstErr = webhook.WebHookToFeishu(msgData2Feishu)

	if dstErr != nil {
		log.Errorf("[NotifyRobot]%v", dstErr.Error())
	}

	return dstErr

}

var (
	gNotifyRobot NotifyRobot
)

func AccessNotify(c *gin.Context) {
	gNotifyRobot.Notify(c)
	c.Next()
}
