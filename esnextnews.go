package main

import (
	"context"
	"esnextnews/parser"
	"fmt"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/spf13/viper"
)

type sesEvent struct {
	Name string `json:"name"`
}

func getTelegremMessage(text string) string {
	return fmt.Sprintf("chat_id=%s&parse_mode=HTML&disable_web_page_preview=true&text=%s",
		viper.GetString("chat_id"),
		url.QueryEscape(text))
}

func getTelegramURL() string {
	return fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", viper.GetString("bot_token"))
}

func handleRequest(ctx context.Context, sesEvent events.SimpleEmailEvent) error {
	if len(sesEvent.Records) < 1 {
		return fmt.Errorf("SES event has no records")
	}

	messageID := sesEvent.Records[0].SES.Mail.MessageID
	svc := s3.New(session.New())
	o, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("aws_bucket")),
		Key:    aws.String(messageID),
	})
	if err != nil {
		return fmt.Errorf("Cannot load data from s3: %s", err)
	}

	md, err := parser.Parse(o.Body)
	if err != nil {
		return fmt.Errorf("Cannot parse loaded data: %s", err)
	}

	_, err = http.Get(getTelegramURL() + getTelegremMessage(md))
	if err != nil {
		return fmt.Errorf("Cannot call Telegram API: %s", err)
	}

	return nil
}

func main() {
	lambda.Start(handleRequest)
}

func init() {
	viper.SetEnvPrefix("enn")
	viper.BindEnv("bot_token")
	viper.BindEnv("chat_id")
	viper.BindEnv("aws_bucket")
}
