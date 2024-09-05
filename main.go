package main

import (
	"flag"
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

	log.Println("The checking of projects has started.")
	for {
		projects, err := services.GetProjects(glc, config)
		if err != nil {
			log.Printf("Error retrieving projects: %v\n", err)
			return
		}

		for _, project := range projects {
			log.Printf("Processing project \"%s\"\n", project.Name)
			services.CheckBranchesAndCreateMR(glc, nc, config, project)
			log.Printf("Finished with project \"%s\"\n", project.Name)
		}

		log.Println("Finished checking branches and processing MRs. Waiting for the next run...")

		time.Sleep(time.Duration(config.CheckIntervalInHours) * time.Hour)
	}
}
