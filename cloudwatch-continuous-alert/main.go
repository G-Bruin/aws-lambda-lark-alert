package main

import (
	"bytes"
	"cloudwatch-continuous-alert/vo"
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"net/http"
	"os"
)

// 处理函数
func handler(ctx context.Context, event vo.AlarmEvent) {

	var reasonData vo.ReasonData
	if err := json.Unmarshal([]byte(event.Detail.State.ReasonData), &reasonData); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	metricName := ""
	var subject map[string]string
	if len(event.Detail.Configuration.Metrics) > 0 {
		metricName = event.Detail.Configuration.Metrics[0].MetricStat.Metric.Name
		subject = event.Detail.Configuration.Metrics[0].MetricStat.Metric.Dimensions
	}

	text := fmt.Sprintf(`
告警策略: %s
告警对象: %v
触发时间: %s
附加通知内容: %s
区域:  %s
度量名称:  %s
阈值: %f`, event.Detail.AlarmName, subject, event.Time, event.Detail.State.Reason, event.Region,
		metricName, reasonData.Threshold)
	PushFeishu(text, event.Detail.AlarmName)

}

func main() {
	lambda.Start(handler)
}

type FeishuMessageCard struct {
	MsgType string `json:"msg_type"`
	Card    struct {
		Elements []Elements `json:"elements"`
		Header   struct {
			Title struct {
				Content string `json:"content"`
				Tag     string `json:"tag"`
			} `json:"title"`
		} `json:"header"`
	} `json:"card"`
}

type Elements struct {
	Tag  string `json:"tag"`
	Text struct {
		Content string `json:"content"`
		Tag     string `json:"tag"`
	} `json:"text,omitempty"`
}

func PushFeishu(content string, title string) {
	secret := os.Getenv("WEBHOOK_KEY")
	webhookURL := "https://open.feishu.cn/open-apis/bot/v2/hook/" + secret

	message := FeishuMessageCard{
		MsgType: "interactive",
	}
	message.Card.Header.Title.Content = title
	message.Card.Header.Title.Tag = "plain_text"

	var arr []Elements
	elements := Elements{
		Tag: "div",
		Text: struct {
			Content string `json:"content"`
			Tag     string `json:"tag"`
		}{
			Content: content,
			Tag:     "lark_md",
		},
	}
	message.Card.Elements = append(arr, elements)

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
