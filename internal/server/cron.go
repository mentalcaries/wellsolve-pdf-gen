package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3"
)

func initializeNotificationCron() {

	mainAppUrl := os.Getenv("MAIN_APP_URL")
	if mainAppUrl == "" {
		log.Fatal("main app url must be set for cron jobs")
	}
	cronSecret := os.Getenv("CRON_SECRET")
	if cronSecret == "" {
		log.Fatal("cron secret must be set")
	}

	c := cron.New()

	c.AddFunc("0 6 1 * *", func() {
		currentMonth := time.Now().Month().String()
		fmt.Println("Running Well Check Notifications Job for", currentMonth)

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/cron/well-reminders", mainAppUrl), nil)
		if err != nil {
			fmt.Printf("couldn't create request. Error: %v\n", err.Error())
			return

		}
		req.Header.Set("Authorization", "Bearer "+cronSecret)

		client := &http.Client{Timeout: 60 * time.Second}
		res, err := client.Do(req)

		if err != nil {
			fmt.Printf("unable to post request to cron endpoint. Error: %v\n", err.Error())
			return
		}
		if res.StatusCode > 299 {
			fmt.Printf("error status code: %v", res.StatusCode)
			return
		}

		defer res.Body.Close()
		fmt.Println("Notifications triggered succesfully with code", res.StatusCode)
	})

	c.Start()
}
