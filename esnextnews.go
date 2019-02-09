package main

import (
	"context"
	"esnextnews/parser"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/credentials"
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

	log.Println("Creating AWS session...")
	region := viper.GetString("aws_region")
	keyID := viper.GetString("aws_key")
	secret := viper.GetString("aws_secret")
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(keyID, secret, ""),
	})
	if err != nil {
		return fmt.Errorf("Cannot create AWS session: %s", err)
	}
	log.Println("AWS session created")

	log.Println("Connecting to S3...")
	svc := s3.New(sess)
	o, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(viper.GetString("aws_bucket")),
		Key:    aws.String(sesEvent.Records[0].SES.Mail.MessageID),
	})
	if err != nil {
		return fmt.Errorf("Cannot load data from s3: %s", err)
	}
	log.Println("S3 data loaded")

	log.Println("Parsing email...")
	md, err := parser.Parse(o.Body)
	if err != nil {
		return fmt.Errorf("Cannot parse loaded data: %s", err)
	}
	log.Println("Email parsed")

	url := getTelegramURL() + "?" + getTelegremMessage(md)
	log.Printf("Calling Telegram API with %s...\n", url)
	_, err = http.Get(url)
	if err != nil {
		return fmt.Errorf("Cannot call Telegram API: %s", err)
	}
	log.Println("Telegram API called")

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
	viper.BindEnv("aws_key")
	viper.BindEnv("aws_secret")
	viper.BindEnv("aws_region")
}
