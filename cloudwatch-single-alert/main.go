package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

const (
	GreaterThanOrEqualToThreshold = "GreaterThanOrEqualToThreshold"
	GreaterThanThreshold          = "GreaterThanThreshold"
	LessThanOrEqualToThreshold    = "LessThanOrEqualToThreshold"
	LessThanThreshold             = "LessThanThreshold"
)

func symbolic(str string) string {
	switch str {
	case GreaterThanOrEqualToThreshold:
		return ">="
	case GreaterThanThreshold:
		return ">"
	case LessThanOrEqualToThreshold:
		return "<="
	case LessThanThreshold:
		return "<"
	default:
		return str
	}
}

// 处理函数
func handler(ctx context.Context, snsEvent events.SNSEvent) error {

	appName := os.Getenv("APP_NAME")
	fmt.Println(appName)
	fmt.Println(snsEvent)

	for _, record := range snsEvent.Records {
		snsRecord := record.SNS
		fmt.Println(snsRecord.Message)

		detail := &events.CloudWatchAlarmSNSPayload{}
		err := json.Unmarshal([]byte(snsRecord.Message), detail)
		if err != nil {
			text := fmt.Sprintf("%s Error unmarshalling SNS message: %v", appName, err)
			fmt.Println(text)
			PushFeishu(text)
			return err
		}

		subject := ""
		if len(detail.Trigger.Dimensions) > 0 {
			subject = detail.Trigger.Dimensions[0].Value
		}

		text := fmt.Sprintf(`
告警策略: %s
告警对象: %s
触发时间: %s
附加通知内容: %s
区域:  %s
度量名称:  %s
阈值: %s %f`, detail.AlarmName, subject, detail.StateChangeTime, detail.NewStateReason, detail.Region,
			detail.Trigger.MetricName, symbolic(detail.Trigger.ComparisonOperator), detail.Trigger.Threshold)
		PushFeishu(text)
	}
	return nil
}

func main() {
	lambda.Start(handler)
}

type FeishuMessage struct {
	MsgType string `json:"msg_type"`
	Content struct {
		Text string `json:"text"`
	} `json:"content"`
}

func PushFeishu(content string) {
	secret := os.Getenv("WEBHOOK_KEY")
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/" + secret

	message := FeishuMessage{
		MsgType: "text",
	}
	message.Content.Text = content

	messageBytes, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshaling message: ", err)
		return
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageBytes))
	if err != nil {
		fmt.Println("Error sending request: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Non-OK HTTP status: ", resp.StatusCode)
		return
	}
}
