package webhook

import (
	"fmt"
	"github.com/OpenListTeam/OpenList/v4/pkg/webhook/feishu"
	"os"
)

func WebHookIsSet() bool {

	dstRet := false

	webhookUrl := os.Getenv("NOTIFY_WEBHOOK_FEISHU")
	if len(webhookUrl) > 0 {
		dstRet = true
		return dstRet
	}

	return dstRet

}

func WebHookToFeishu(msgData string) error {

	var dstErr error

	webhookUrl := os.Getenv("NOTIFY_WEBHOOK_FEISHU")
	if len(webhookUrl) < 1 {
		return dstErr
	}

	client := feishu.NewClient(webhookUrl)

	sendErr := client.SendText(msgData)

	if sendErr != nil {
		dstErr = fmt.Errorf("FeishuClient SendText error:%v", sendErr.Error())
	}

	return dstErr

}
