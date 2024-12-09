package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/cyinnove/subfalcon/pkg/takeover"
)

type DiscordWebhook struct {
	Content string  `json:"content,omitempty"`
	Embeds  []Embed `json:"embeds,omitempty"`
}

type Embed struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Color       int     `json:"color"`
	Fields      []Field `json:"fields,omitempty"`
}

type Field struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline"`
}

const (
	MaxMessageLength = 2000
	ChunkSize        = 20 // number of subdomains per message
)

func SendTakeoverResults(webhookURL string, results []*takeover.SubdomainInfo) error {
	if len(results) == 0 {
		return nil
	}

	// If results are too large, save to file and send as attachment
	if len(results) > 50 {
		return sendTakeoverResultsAsFile(webhookURL, results)
	}

	// Split results into chunks
	chunks := chunkTakeoverResults(results, ChunkSize)

	for i, chunk := range chunks {
		webhook := DiscordWebhook{
			Embeds: []Embed{
				{
					Title: fmt.Sprintf("🔍 Subdomain Takeover Results (Part %d/%d)", i+1, len(chunks)),
					Color: 0xff0000,
					Fields: []Field{
						{
							Name:  "Results",
							Value: formatTakeoverChunk(chunk),
						},
					},
				},
			},
		}

		if err := sendWebhook(webhookURL, webhook); err != nil {
			return err
		}

		// Avoid Discord rate limits
		time.Sleep(1 * time.Second)
	}

	return nil
}

func formatTakeoverChunk(results []*takeover.SubdomainInfo) string {
	var sb strings.Builder
	for _, info := range results {
		if info.CNAME != "" {
			status := "Active"
			if info.IsVulnerable {
				status = "⚠️ VULNERABLE"
			}
			sb.WriteString(fmt.Sprintf("• %s\n  ↳ CNAME: %s\n  ↳ Status: %s\n",
				info.Subdomain, info.CNAME, status))
		}
	}
	return sb.String()
}

func sendTakeoverResultsAsFile(webhookURL string, results []*takeover.SubdomainInfo) error {
	// Create temporary file
	tmpfile, err := os.CreateTemp("", "takeover-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	// Write results to file
	for _, info := range results {
		if info.CNAME != "" {
			fmt.Fprintf(tmpfile, "Subdomain: %s\nCNAME: %s\nStatus: %s\nVulnerable: %v\n\n",
				info.Subdomain, info.CNAME, info.Status, info.IsVulnerable)
		}
	}
	tmpfile.Close()

	// Send file to Discord
	return sendFileToDiscord(webhookURL, tmpfile.Name(), "takeover-results.txt")
}

func sendFileToDiscord(webhookURL, filepath, filename string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return err
	}

	writer.Close()

	req, err := http.NewRequest("POST", webhookURL, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("discord webhook returned status code %d", resp.StatusCode)
	}

	return nil
}

func SendNewSubdomains(webhookURL string, newSubdomains []string) error {
	if len(newSubdomains) == 0 {
		return nil
	}

	// If too many subdomains, send as file
	if len(newSubdomains) > 100 {
		return sendSubdomainsAsFile(webhookURL, newSubdomains)
	}

	// Split into chunks
	chunks := chunkSlice(newSubdomains, ChunkSize)

	for i, chunk := range chunks {
		webhook := DiscordWebhook{
			Embeds: []Embed{
				{
					Title:       fmt.Sprintf("🆕 New Subdomains (Part %d/%d)", i+1, len(chunks)),
					Color:       0x3498db,
					Description: "```\n" + strings.Join(chunk, "\n") + "\n```",
				},
			},
		}

		if err := sendWebhook(webhookURL, webhook); err != nil {
			return err
		}

		// Avoid Discord rate limits
		time.Sleep(1 * time.Second)
	}

	return nil
}

func sendSubdomainsAsFile(webhookURL string, subdomains []string) error {
	tmpfile, err := os.CreateTemp("", "subdomains-*.txt")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	for _, subdomain := range subdomains {
		fmt.Fprintln(tmpfile, subdomain)
	}
	tmpfile.Close()

	return sendFileToDiscord(webhookURL, tmpfile.Name(), "new-subdomains.txt")
}

func sendWebhook(webhookURL string, webhook DiscordWebhook) error {
	jsonData, err := json.Marshal(webhook)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return fmt.Errorf("discord webhook returned status code %d", resp.StatusCode)
	}

	return nil
}

func chunkSlice(slice []string, chunkSize int) [][]string {
	var chunks [][]string
	for i := 0; i < len(slice); i += chunkSize {
		end := i + chunkSize
		if end > len(slice) {
			end = len(slice)
		}
		chunks = append(chunks, slice[i:end])
	}
	return chunks
}

func chunkTakeoverResults(results []*takeover.SubdomainInfo, chunkSize int) [][]*takeover.SubdomainInfo {
	var chunks [][]*takeover.SubdomainInfo
	for i := 0; i < len(results); i += chunkSize {
		end := i + chunkSize
		if end > len(results) {
			end = len(results)
		}
		chunks = append(chunks, results[i:end])
	}
	return chunks
}
