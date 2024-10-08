package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"branch-synchronizer/env"
	"branch-synchronizer/internal/services"
	"branch-synchronizer/pkg/utils"

	"github.com/xanzy/go-gitlab"
)

func main() {
	logFile, err := utils.NewLogger()
	if err != nil {
		log.Printf("Error configure logger: %v\n", err)
		return
	}
	defer logFile.Close()

	configPath := flag.String("config", "", "Path to configuration file")
	flag.Parse()

	log.Println("Loading configuration...")
	config, err := env.NewConfig(*configPath)
	if err != nil {
		log.Printf("Error creating configuration instance: %v\n", err)
		return
	}

	log.Println("Creating GitLab client...")
	glc, err := gitlab.NewClient(config.GitLab.Token, gitlab.WithBaseURL(config.GitLab.URL))
	if err != nil {
		log.Printf("Error creating gitlab client: %v\n", err)
		return
	}

	log.Println("Creating Notifier(s) Bot client...")
	nc, err := services.NewNotifierClient(config)
	if err != nil {
		log.Printf("Error creating Notifier(s) bot client: %v\n", err)
		return
	}

	if config.CheckIntervalInHours <= 0 {
		log.Println("The checking of projects has started.")
		checkProjects(glc, nc, config)
	} else {
		for {
			log.Println("The checking of projects has started.")
			checkProjects(glc, nc, config)

			log.Println("Waiting for the next run...")
			time.Sleep(time.Duration(config.CheckIntervalInHours) * time.Hour)
		}
	}
}

func checkProjects(glc *gitlab.Client, nc *services.NotifierClient, config *env.Config) {
	projects, err := services.GetProjects(glc, config)
	if err != nil {
		log.Printf("Error retrieving projects: %v\n", err)
		log.Println("Without projects it's impossible to do anything further.")
		return
	}

	var mrCounter uint
	for _, project := range projects {
		log.Printf("Processing project \"%s\"\n", project.Name)
		services.CheckBranchesAndCreateMR(glc, nc, config, project, &mrCounter)
		log.Printf("Finished with project \"%s\"\n", project.Name)

		time.Sleep(10 * time.Second)
	}

	if mrCounter == 0 {
		message := fmt.Sprintf("No MRs have been created in the current run.\nGoing back to sleep 😴")
		err = nc.SendNotification(message)
		if err != nil {
			log.Printf("Error sending notification about absence of MRs creation: %v\n", err)
		}
	}

	log.Println("Finished checking branches and processing MRs.")
}
