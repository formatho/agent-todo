package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/formatho/agent-todo/cli/client"
	"github.com/spf13/cobra"
)

type Agent struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Role        string `json:"role"`
	Enabled     bool   `json:"enabled"`
	APIKey      string `json:"api_key"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var (
	agentEnabled bool
	agentRole    string
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Agent management commands",
	Long:  "Commands for managing AI agents",
}

var createAgentCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description, _ := cmd.Flags().GetString("description")
		role, _ := cmd.Flags().GetString("role")

		// Default to regular role if not specified
		if role == "" {
			role = "regular"
		}

		// Validate role
		if role != "regular" && role != "supervisor" && role != "admin" {
			return fmt.Errorf("invalid role: %s (must be regular, supervisor, or admin)", role)
		}

		c := client.New()
		req := map[string]interface{}{
			"name": name,
			"role": role,
		}
		if description != "" {
			req["description"] = description
		}

		resp, err := c.Post("/agents", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to create agent: %s", string(body))
		}

		var agent Agent
		if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Agent created: %s (ID: %s)\n", agent.Name, agent.ID)
		fmt.Printf("Role: %s\n", agent.Role)
		fmt.Printf("API Key: %s\n", agent.APIKey)
		fmt.Println("\n⚠ Save this API key securely. It won't be shown again.")
		return nil
	},
}

var listAgentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all agents",
	RunE: func(cmd *cobra.Command, args []string) error {
		c := client.New()
		resp, err := c.Get("/agents", true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to list agents: %s", string(body))
		}

		var agents []Agent
		if err := json.NewDecoder(resp.Body).Decode(&agents); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		if len(agents) == 0 {
			fmt.Println("No agents found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tROLE\tENABLED")
		for _, a := range agents {
			enabled := "✓"
			if !a.Enabled {
				enabled = "✗"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", a.ID, a.Name, a.Role, enabled)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d agent(s)\n", len(agents))
		return nil
	},
}

var getAgentCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get agent details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Get("/agents/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get agent: %s", string(body))
		}

		var agent Agent
		if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("ID:          %s\n", agent.ID)
		fmt.Printf("Name:        %s\n", agent.Name)
		fmt.Printf("Description: %s\n", agent.Description)
		fmt.Printf("Role:        %s\n", agent.Role)
		fmt.Printf("Enabled:     %v\n", agent.Enabled)
		fmt.Printf("API Key:     %s\n", agent.APIKey)
		fmt.Printf("Created:     %s\n", agent.CreatedAt)
		fmt.Printf("Updated:     %s\n", agent.UpdatedAt)
		return nil
	},
}

var updateAgentCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		role, _ := cmd.Flags().GetString("role")
		enabled, _ := cmd.Flags().GetBool("enabled")

		req := make(map[string]interface{})
		if name != "" {
			req["name"] = name
		}
		if description != "" {
			req["description"] = description
		}
		if role != "" {
			// Validate role
			if role != "regular" && role != "supervisor" && role != "admin" {
				return fmt.Errorf("invalid role: %s (must be regular, supervisor, or admin)", role)
			}
			req["role"] = role
		}
		if cmd.Flags().Changed("enabled") {
			req["enabled"] = enabled
		}

		if len(req) == 0 {
			return fmt.Errorf("at least one field must be specified")
		}

		c := client.New()
		resp, err := c.Patch("/agents/"+id, req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update agent: %s", string(body))
		}

		var agent Agent
		if err := json.NewDecoder(resp.Body).Decode(&agent); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Agent updated: %s (Role: %s)\n", agent.Name, agent.Role)
		return nil
	},
}

var deleteAgentCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Delete("/agents/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete agent: %s", string(body))
		}

		fmt.Println("✓ Agent deleted")
		return nil
	},
}

var agentListTasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List tasks assigned to the authenticated agent",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")

		path := "/agent/tasks"
		if status != "" {
			path += "?status=" + status
		}

		c := client.New()
		resp, err := c.Get(path, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to list tasks: %s", string(body))
		}

		var tasks []Task
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY")
		for _, t := range tasks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", t.ID, t.Title, t.Status, t.Priority)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d task(s)\n", len(tasks))
		return nil
	},
}

var agentGetTaskCmd = &cobra.Command{
	Use:   "get-task <id>",
	Short: "Get details of a task assigned to the agent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Get("/agent/tasks/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get task: %s", string(body))
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("ID:          %s\n", task.ID)
		fmt.Printf("Title:       %s\n", task.Title)
		fmt.Printf("Description: %s\n", task.Description)
		fmt.Printf("Status:      %s\n", task.Status)
		fmt.Printf("Priority:    %s\n", task.Priority)
		fmt.Printf("Created:     %s\n", task.CreatedAt)
		fmt.Printf("Updated:     %s\n", task.UpdatedAt)
		return nil
	},
}

var agentUpdateTaskStatusCmd = &cobra.Command{
	Use:   "update-status <id> <status>",
	Short: "Update the status of a task assigned to the agent",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		status := args[1]

		// Validate status
		if status != "pending" && status != "in_progress" && status != "completed" && status != "failed" {
			return fmt.Errorf("invalid status: %s (must be pending, in_progress, completed, or failed)", status)
		}

		c := client.New()
		req := map[string]string{
			"status": status,
		}

		resp, err := c.Patch("/agent/tasks/"+id+"/status", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update task status: %s", string(body))
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Task status updated: %s -> %s\n", task.Title, task.Status)
		return nil
	},
}

var agentAddCommentCmd = &cobra.Command{
	Use:   "add-comment <task-id> <comment>",
	Short: "Add a comment to a task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		content := args[1]

		c := client.New()
		req := map[string]string{
			"content": content,
		}

		resp, err := c.Post("/agent/tasks/"+taskID+"/comments", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to add comment: %s", string(body))
		}

		fmt.Println("✓ Comment added successfully")
		return nil
	},
}

var agentListCommentsCmd = &cobra.Command{
	Use:   "comments <task-id>",
	Short: "List comments on a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		c := client.New()
		resp, err := c.Get("/agent/tasks/"+taskID+"/comments", true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to list comments: %s", string(body))
		}

		var comments []Comment
		if err := json.NewDecoder(resp.Body).Decode(&comments); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		if len(comments) == 0 {
			fmt.Println("No comments found")
			return nil
		}

		for _, comment := range comments {
			fmt.Printf("[%s] %s: %s\n\n", comment.CreatedAt, comment.AuthorName, comment.Content)
		}

		fmt.Printf("Total: %d comment(s)\n", len(comments))
		return nil
	},
}

var describeAgentCmd = &cobra.Command{
	Use:   "describe [name-or-id]",
	Short: "Get agent description in markdown format",
	Long:  "Output a detailed agent description in markdown format, including capabilities, access, and configuration. If no name/id provided, describes the current authenticated agent.",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFile, _ := cmd.Flags().GetString("output")

		c := client.New()

		// Get current agent's tasks to verify auth and get agent info
		resp, err := c.Get("/agent/tasks?limit=1", true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to authenticate: %s", string(body))
		}

		// Parse response to extract agent info from any task
		var tasks []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		// Build agent info from config and task data
		agent := &Agent{
			ID:          "current",
			Name:        "current-agent",
			Description: "The agent authenticated via the configured API key",
			Role:        "regular",
			Enabled:     true,
		}

		// If we got tasks, extract more info
		if len(tasks) > 0 {
			task := tasks[0]
			if assignedAgent, ok := task["assigned_agent"].(map[string]interface{}); ok {
				if id, ok := assignedAgent["id"].(string); ok {
					agent.ID = id
				}
				if name, ok := assignedAgent["name"].(string); ok {
					agent.Name = name
				}
				if role, ok := assignedAgent["role"].(string); ok {
					agent.Role = role
				}
				if desc, ok := assignedAgent["description"].(string); ok {
					agent.Description = desc
				}
			}
		}

		outputMarkdown(agent, outputFile)
		return nil
	},
}

func outputMarkdown(agent *Agent, outputFile string) {
	markdown := fmt.Sprintf(`# Agent: %s

## Overview

| Property | Value |
|----------|-------|
| ID | %s |
| Name | %s |
| Role | %s |
| Status | %s |
| Created | %s |

## Description

%s

## Capabilities

### Role Permissions

`, agent.Name, agent.ID, agent.Name, agent.Role, enabledStatus(agent.Enabled), agent.CreatedAt, agent.Description)

	// Add role-based capabilities
	switch agent.Role {
	case "admin":
		markdown += `- **Full system access** - Can perform any action
- **Agent management** - Create, update, delete agents
- **Task management** - Full CRUD on all tasks
- **User management** - Manage users and permissions
- **System configuration** - Modify system settings
`
	case "supervisor":
		markdown += `- **Task oversight** - View and update any task
- **Agent creation** - Create new agents
- **Task assignment** - Assign tasks to agents
- **Reporting** - Access to reports and metrics
`
	default:
		markdown += `- **Own tasks only** - Can only update tasks assigned to them
- **Task creation** - Can create new tasks
- **Comments** - Can add comments to assigned tasks
- **Status updates** - Can update status of assigned tasks
`
	}

	markdown += fmt.Sprintf(`
## API Access

This agent can authenticate using the X-API-KEY header:

` + "```" + `bash
curl -H "X-API-Key: <api-key>" https://todo.formatho.com/api/agent/tasks
` + "```" + `

## Available Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| /agent/tasks | GET | List assigned tasks |
| /agent/tasks | POST | Create a new task |
| /agent/tasks/:id | GET | Get task details |
| /agent/tasks/:id/status | PATCH | Update task status |
| /agent/tasks/:id/comments | GET | List task comments |
| /agent/tasks/:id/comments | POST | Add a comment |

## Usage Guidelines

1. **Authentication**: Always include the API key in requests
2. **Task Updates**: Update status promptly when working on tasks
3. **Comments**: Document progress and blockers in comments
4. **Priority**: Work on high-priority tasks first

---
*Generated by agent-todo CLI*
`, )

	if outputFile != "" {
		err := os.WriteFile(outputFile, []byte(markdown), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %v\n", err)
			fmt.Println(markdown)
			return
		}
		fmt.Printf("✓ Agent description written to: %s\n", outputFile)
	} else {
		fmt.Println(markdown)
	}
}

func enabledStatus(enabled bool) string {
	if enabled {
		return "✓ Active"
	}
	return "✗ Disabled"
}

func init() {
	rootCmd.AddCommand(agentCmd)

	// Create agent
	createAgentCmd.Flags().StringP("description", "d", "", "Agent description")
	createAgentCmd.Flags().StringP("role", "r", "", "Agent role (regular, supervisor, admin)")

	// Update agent
	updateAgentCmd.Flags().StringP("name", "n", "", "Agent name")
	updateAgentCmd.Flags().StringP("description", "d", "", "Agent description")
	updateAgentCmd.Flags().StringP("role", "r", "", "Agent role (regular, supervisor, admin)")
	updateAgentCmd.Flags().BoolVarP(&agentEnabled, "enabled", "e", false, "Enable/disable agent")

	// Agent tasks
	agentListTasksCmd.Flags().StringP("status", "s", "", "Filter by status")

	// Describe agent
	describeAgentCmd.Flags().StringP("output", "o", "", "Output file (prints to stdout if not specified)")

	agentCmd.AddCommand(createAgentCmd)
	agentCmd.AddCommand(listAgentsCmd)
	agentCmd.AddCommand(getAgentCmd)
	agentCmd.AddCommand(updateAgentCmd)
	agentCmd.AddCommand(deleteAgentCmd)
	agentCmd.AddCommand(describeAgentCmd)
	agentCmd.AddCommand(agentListTasksCmd)
	agentCmd.AddCommand(agentGetTaskCmd)
	agentCmd.AddCommand(agentUpdateTaskStatusCmd)
	agentCmd.AddCommand(agentAddCommentCmd)
	agentCmd.AddCommand(agentListCommentsCmd)
}
