package main

import (
	"context"
	"log"
	"os"
	"time"

	"google.golang.org/genai"
	"gorm.io/gorm"
)

var GeminiClient *genai.Client
var DB *gorm.DB

const (
	GModel     = "gemma-3-27b-it"
	GBigModel  = "gemini-3-flash-preview"
	GBigModel2 = "gemini-2.0-flash"
)

func main() {
	log.Println("Hello World! We have started!")

	GeminiKey := os.Getenv("WORLDGAME_GEMINI")

	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey:  GeminiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		log.Fatalf("Could not get GEMINI to work... %+v\n", err)
	}
	GeminiClient = client

	InitDB("localhost", "postgres", os.Getenv("POSTGRES_PASSWORD"))

	go SetupWebServer()

	//GetTagesschauThemes()
	//set timers for time-based tasks
	Timer()
}

func Timer() {
	loc, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		log.Fatal("No timezone info / Could not fetch timezone info!!!")
	}
	for {
		now := time.Now()

		nextRun := time.Date(now.Year(), now.Month(), now.Day(), 7, 0, 0, 0, loc)

		if nextRun.Before(now) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		timer := time.NewTimer(nextRun.Sub(now))
		log.Printf("Scheduled next run for %s\n", nextRun)

		<-timer.C

		// TODO: replace with whatever will come up
		// go GetTagesschauThemes()
	}
}
