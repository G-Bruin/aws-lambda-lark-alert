# aws-lambda-lark-alert

## 简介
- cloudwatch-single-alert
    - 使用 Lambda、SNS、Cloudwatch 指标进行飞书的告警监控
- cloudwatch-continuous-alert
    - 由于 aws 指标监控的限制，仅仅变动的时候才会进行push，无法做到连续时间的持续监控，故进行升级
    - 使用 Lambda、Step Functions、Amazon EventBridge、Cloudwatch 指标进行飞书的告警监控
    - [持续报警参考文档](https://aws.amazon.com/cn/blogs/china/use-aws-step-functions-to-implement-continuous-amazon-cloudwatch-alarms/)
- 需要给role配置 CheckAlarmDescribeAlarms 权限
- 运行时配置
  ![image](https://github.com/user-attachments/assets/f939b419-6d40-486b-b9bb-11c35730b200)

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Action": [
                "cloudwatch:DescribeAlarms*"
            ],
            "Effect": "Allow",
            "Resource": "*"
        }
    ]
}
```

## 通用使用方法

### 编译 GO 文件

```bash
go mod tidy
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o bootstrap -tags lambda.norpc main.go
zip myFunction.zip bootstrap
```


### 上传和更新 Lambda 函数
```bash

aws lambda create-function --function-name test-go-sns-event --profile sg  --runtime provided.al2023 --handler bootstrap --architectures arm64 --role arn:aws:iam::{accountId}:role/lambda-ex --zip-file fileb://myFunction.zip

aws lambda update-function-code --function-name test-go-sns-event --zip-file fileb://myFunction.zip
```

### 设置环境变量

```bash
- lambda 设置名为 `WEBHOOK_KEY` 的环境变量，其值为自己的飞书机器人 webhook url 的最后哈希串。
- 设置飞书 webhook 为：`https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx-xxxx-xxx-xxx-xxx`。
- `WEBHOOK_KEY` 需要设置一个值为 `xxxxx-xxxx-xxx-xxx-xxx` 的环境变量 
```

### aws-cli 配置环境变量
```
vim ~/.aws/credentials

[default]
aws_access_key_id = xx
aws_secret_access_key = xx
[sg]
aws_access_key_id = xx
aws_secret_access_key = xx

```

```
vim ~/.aws/config

[default]
region = us-east-1

[profile sg]
region = ap-southeast-1

```
