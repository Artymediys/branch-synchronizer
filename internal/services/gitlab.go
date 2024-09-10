package services

import (
	"fmt"
	"log"
	"strings"

	"branch-synchronizer/env"

	"github.com/xanzy/go-gitlab"
)

func CheckBranchesAndCreateMR(
	glc *gitlab.Client,
	notifier *NotifierClient,
	config *env.Config,
	project *gitlab.Project,
	mrCounter *uint,
) {
	for _, pair := range config.BranchPairs {
		branches := strings.Split(pair, "->")
		if len(branches) != 2 {
			log.Printf("Invalid branch pair: %s\n", pair)
			continue
		}

		sourceBranch := strings.TrimSpace(branches[0])
		targetBranch := strings.TrimSpace(branches[1])

		baseCompare, _, err := glc.Repositories.Compare(project.ID, &gitlab.CompareOptions{
			From: &sourceBranch,
			To:   &targetBranch,
		})
		if err != nil {
			log.Printf("Error base comparing branches for project %s: %v\n", project.Name, err)
			continue
		}

		isStraight := true
		straightCompare, _, err := glc.Repositories.Compare(project.ID, &gitlab.CompareOptions{
			From:     &sourceBranch,
			To:       &targetBranch,
			Straight: &isStraight,
		})
		if err != nil {
			log.Printf("Error straight comparing branches for project %s: %v\n", project.Name, err)
			continue
		}

		if len(baseCompare.Diffs) > 0 || len(straightCompare.Diffs) > 0 {
			mrURL, err := createMR(glc, project.ID, sourceBranch, targetBranch)
			if err != nil {
				log.Printf("Error creating MR for project %s: %v\n", project.Name, err)
				continue
			}

			log.Printf("MR created: %s\n", mrURL)
			message := fmt.Sprintf(
				"> Created MR for project **\"%s\"**\nBranches: `%s -> %s`\nMerge Request: [link](%s)",
				project.Name, sourceBranch, targetBranch, mrURL,
			)
			*mrCounter++

			err = notifier.SendNotification(message)
			if err != nil {
				log.Printf("Error sending notification for project \"%s\": %v\n", project.Name, err)
			}
		}
	}
}

func GetProjects(glc *gitlab.Client, config *env.Config) ([]*gitlab.Project, error) {
	allProjects := make([]*gitlab.Project, 0, 32)

	groupIDs, err := getSubgroupIDs(glc, config.GitLab.GroupID)
	if err != nil {
		return nil, err
	}
	groupIDs = append(groupIDs, config.GitLab.GroupID)

	for _, groupID := range groupIDs {
		options := &gitlab.ListGroupProjectsOptions{
			ListOptions: gitlab.ListOptions{
				PerPage: 50,
				Page:    1,
			},
		}

		for {
			projects, resp, err := glc.Groups.ListGroupProjects(groupID, options)
			if err != nil {
				return nil, err
			}
			allProjects = append(allProjects, projects...)

			if resp.CurrentPage >= resp.TotalPages {
				break
			}

			options.Page = resp.NextPage
		}
	}

	filteredProjects := make([]*gitlab.Project, 0, 4)

	for _, project := range allProjects {
		projectSysName := getProjectPath(project)

		switch {
		case len(config.WhitelistProjects) > 0:
			if contains(config.WhitelistProjects, projectSysName) {
				filteredProjects = append(filteredProjects, project)
			}
		case len(config.BlacklistProjects) > 0:
			if !contains(config.BlacklistProjects, projectSysName) {
				filteredProjects = append(filteredProjects, project)
			}
		default:
			filteredProjects = append(filteredProjects, project)
		}
	}

	return filteredProjects, nil
}

func getSubgroupIDs(glc *gitlab.Client, mainGroupID int) ([]int, error) {
	opt := &gitlab.ListSubGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 50,
			Page:    1,
		},
	}

	allSubgroupIDs := make([]int, 0, 4)

	for {
		subgroups, resp, err := glc.Groups.ListSubGroups(mainGroupID, opt)
		if err != nil {
			return nil, err
		}

		for _, subgroup := range subgroups {
			allSubgroupIDs = append(allSubgroupIDs, subgroup.ID)

			subSubGroupIDs, err := getSubgroupIDs(glc, subgroup.ID)
			if err != nil {
				return nil, err
			}

			allSubgroupIDs = append(allSubgroupIDs, subSubGroupIDs...)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allSubgroupIDs, nil
}

func getProjectPath(project *gitlab.Project) string {
	urlParts := strings.Split(project.WebURL, "/")
	return urlParts[len(urlParts)-1]
}

func contains(list []string, item string) bool {
	for _, elem := range list {
		if elem == item {
			return true
		}
	}

	return false
}

func createMR(glc *gitlab.Client, projectID int, sourceBranch, targetBranch string) (string, error) {
	mrTitle := fmt.Sprintf("Merge %s into %s", sourceBranch, targetBranch)
	mrOptions := &gitlab.CreateMergeRequestOptions{
		SourceBranch: &sourceBranch,
		TargetBranch: &targetBranch,
		Title:        &mrTitle,
	}

	mr, _, err := glc.MergeRequests.CreateMergeRequest(projectID, mrOptions)
	if err != nil {
		return "", err
	}

	return mr.WebURL, nil
}
