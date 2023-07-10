package main

import (
	"fmt"
	"log"

	cloudranges "github.com/dvvedz/cloudranges/pkg"
)

func main() {
	var githubRanges, grErr = cloudranges.GithubRanges("v4")
	if grErr != nil {
		log.Fatalf("something went wrong with GithubRanges, err: %v", grErr)
	}

	// var bingbotRanges, brErr = cloudranges.BingbotRanges()
	// if brErr != nil {
	// 	log.Fatalf("something went wrong with Bingbot, err: %v", brErr)
	// }

	var googleCloudRanges, gcrErr = cloudranges.GoogleCloudRanges("v4")
	if gcrErr != nil {
		log.Fatalf("something went wrong with googleCloudRanges, err: %v", gcrErr)
	}

	var googleBotRanges, gbrErr = cloudranges.GoogleBotRanges("v4")
	if gbrErr != nil {
		log.Fatalf("something went wrong with googleBotRanges, err: %v", gbrErr)
	}

	var awsRanges, awsRangesErr = cloudranges.AwsRanges("v4")
	if awsRangesErr != nil {
		log.Fatalf("something went wrong with Awsranges, err: %v", awsRangesErr)
	}

	var oracleCloud, oracleCloudErr = cloudranges.OracleRanges()
	if oracleCloudErr != nil {
		log.Fatalf("something went wrong with OracleCloud, err: %v", oracleCloudErr)
	}

	var linodeRanges, linodeRangesErr = cloudranges.LinodeRanges("v4")
	if linodeRangesErr != nil {
		log.Fatalf("something went wrong with LinodeRanges, err: %v", linodeRangesErr)
	}

	// var telegramRanges, telegramRangesErr = cloudranges.TelegramRanges("all")
	// if telegramRangesErr != nil {
	// 	log.Fatalf("something went wrong with TelegramRanges, err: %v", telegramRangesErr)
	// }

	var combinded = append(githubRanges, googleBotRanges...)
	combinded = append(combinded, googleCloudRanges...)
	combinded = append(combinded, awsRanges...)
	combinded = append(combinded, oracleCloud...)
	combinded = append(combinded, linodeRanges...)

	for _, iprange := range combinded {
		fmt.Println(iprange)
	}
}
