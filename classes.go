package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FinalJobGrow struct {
	Jobs []JobGrowJSON `json:"jobs"`
}

type JobGrow struct {
	JobGrowID string `json:"jobGrowId"`
	JobGrowName string `json:"jobGrowName"`
	JobAdvancement *JobGrow `json:"next"`
}

type JobGrowJSON struct {
	JobGrowID string `json:"jobGrowId" gorm:"primaryKey"`
	JobGrowName string `json:"jobGrowName" gorm:"primaryKey"`
}

type Job struct {
	JobID string `json:"jobId"`
	JobName string `json:"jobName"`
	JobGrow []JobGrow `json:"rows"`
}

type Response struct {
	Data []Job `json:"rows"`
}

const (
	baseURL = "https://api.dfoneople.com/df/jobs?apikey=f7duulCyzOcAdt3jKEoUyNYrBGlFIAhm"
)

func main() {
	// database connection
	dsn := "user=postgres password=st.XojC6ZVee host=db.odnbticprnftyfyhkupr.supabase.co port=5432 dbname=postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}
	db.AutoMigrate(&JobGrowJSON{})

	// api request
	res, err := http.Get(baseURL)

	if err != nil {
		fmt.Printf("Combination %s - API request failed %s\n", "ab", err)
		return
	}

	var response Response
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&response); err != nil {
		fmt.Printf("Combination %s - JSON decoding error: %s\n", "ab", err)
	}

	var finalJobGrow FinalJobGrow

	var baseJobs []Job
	baseJobs = append(baseJobs, response.Data...)
	for _, baseJob := range baseJobs {
		// fmt.Printf("JobID: %s, JobName: %s\n", baseJob.JobID, baseJob.JobName)
		for _, jobAdv := range baseJob.JobGrow {
			// fmt.Printf("  Job Grow: %s\n", jobAdv.JobGrowName)
			curr := &jobAdv
			for curr != nil {
				next := curr.JobAdvancement
				if next == nil {
					fmt.Printf("	Job Grow ID: %s Job Grow Name: %s\n", curr.JobGrowID, curr.JobGrowName)
					jobGrowJson := JobGrowJSON{
						JobGrowID: curr.JobGrowID,
						JobGrowName: curr.JobGrowName,
					}
					finalJobGrow.Jobs = append(finalJobGrow.Jobs, jobGrowJson)

					// insert into db
					result := db.Create(&jobGrowJson)
					if result.Error != nil {
						fmt.Println(result.Error)
					}

				} else {
					// fmt.Printf("	Job Grow: %s\n", curr.JobGrowName)
				}
				curr = next
			}
		}
	}

	file, _ := json.MarshalIndent(finalJobGrow, "", "    ")
	fmt.Println(string(file)) 
	outputFile, err := os.Create("classes.json")
	_, err = outputFile.Write(file)
}