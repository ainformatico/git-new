package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type Response struct {
	Fields struct {
		Summary string `json:"summary"`
	} `json:"fields"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide an issue number")
		os.Exit(1)
	}

	issue := strings.ToLower(os.Args[1])

	if issue == "" {
		fmt.Println("Please provide an issue number")
		os.Exit(1)
	}

	jiraServer := os.Getenv("JIRA_SERVER")
	jiraUser := os.Getenv("JIRA_USER")
	jiraToken := os.Getenv("JIRA_TOKEN")

	if jiraServer == "" || jiraUser == "" || jiraToken == "" {
		fmt.Println("Please set JIRA_SERVER, JIRA_USER, and JIRA_TOKEN")
		os.Exit(1)
	}

	url := fmt.Sprintf("%s/rest/api/2/issue/%s", jiraServer, issue)
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println("Error creating request:", err)
		os.Exit(1)
	}

	req.SetBasicAuth(jiraUser, jiraToken)

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode != 200 {
		fmt.Printf("Error fetching issue: %s, server returned: %s", issue, resp.Status)

		os.Exit(1)
	}

	if err != nil {
		fmt.Println("Error sending request:", err)

		os.Exit(1)
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var data Response

	json.Unmarshal(body, &data)

	summary := strings.ToLower(data.Fields.Summary)
	summary = strings.TrimSpace(summary)
	summary = regexp.MustCompile(`\W+`).ReplaceAllString(summary, "-")
	summary = regexp.MustCompile(`-+`).ReplaceAllString(summary, "-")

	// Remove non-alphabetical characters at the end of the summary
	summary = regexp.MustCompile(`[^a-z]+$`).ReplaceAllString(summary, "")

	// this ensures we don't hit limits in some tools (e.g Netlify)
	summary = summary[:min(len(summary), 69)]

	fmt.Println("Creating branch:", fmt.Sprintf("%s-%s", issue, summary))

	cmd := exec.Command("git", "checkout", "-b", fmt.Sprintf("%s-%s", issue, summary))

	err = cmd.Run()

	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
