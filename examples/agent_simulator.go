package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const apiBase = "http://localhost:3000/api/v1"

func beacon(agentID string) {
	req := map[string]interface{}{
		"agent_id": agentID,
		"hostname": "banco-host-01",
		"username": "cajero01",
		"os":       "Linux",
		"arch":     "x86_64",
		"tags":     []string{"banco", "prod", "pampero"},
		"zone":     "ARG",
		"profile":  "medios",
		"status":   "active",
	}
	jsonBody, _ := json.Marshal(req)
	resp, err := http.Post(apiBase+"/agents/beacon", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Beacon error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Beacon status=%d body=%s", resp.StatusCode, string(body))
}

func pollJobs(agentID string) []map[string]interface{} {
	resp, err := http.Get(apiBase + "/agents/" + agentID + "/jobs")
	if err != nil {
		log.Printf("Poll jobs error: %v", err)
		return nil
	}
	defer resp.Body.Close()
	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Printf("Error decode jobs: %v", err)
		return nil
	}
	if jobs, ok := data["jobs"].([]interface{}); ok {
		result := []map[string]interface{}{}
		for _, j := range jobs {
			if jm, ok := j.(map[string]interface{}); ok {
				result = append(result, jm)
			}
		}
		return result
	}
	return nil
}

func postResult(agentID, jobID string, stdout string, stderr string, exitCode int) {
	req := map[string]interface{}{
		"agent_id":  agentID,
		"exit_code": exitCode,
		"stdout":    stdout,
		"stderr":    stderr,
	}
	jsonBody, _ := json.Marshal(req)
	resp, err := http.Post(apiBase+"/jobs/"+jobID+"/results", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		log.Printf("Post result error: %v", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	log.Printf("Result status=%d body=%s", resp.StatusCode, string(body))
}

func runCommand(cmd string) (string, string, int) {
	out := fmt.Sprintf("executed %s at %s", cmd, time.Now().Format(time.RFC3339))
	if strings.Contains(cmd, "fail") {
		return "", "command failed", 1
	}
	return out, "", 0
}

func main() {
	agentID := "pampero-agent-0001"
	if len(os.Args) > 1 {
		agentID = os.Args[1]
	}

	log.Printf("Starting agent simulator %s", agentID)
	for {
		beacon(agentID)
		jobs := pollJobs(agentID)
		for _, job := range jobs {
			jid := fmt.Sprintf("%v", job["id"])
			cmd := fmt.Sprintf("%v", job["command"])
			out, err, code := runCommand(cmd)
			postResult(agentID, jid, out, err, code)
		}
		delay := time.Duration(20+rand.Intn(30)) * time.Second
		time.Sleep(delay)
	}
}
