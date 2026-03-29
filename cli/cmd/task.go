package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/formatho/agent-todo/cli/client"
	"github.com/spf13/cobra"
)

type Task struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	ProjectID   string `json:"project_id"`
	AgentID     *string `json:"agent_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type Comment struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	UserID     string `json:"user_id"`
	AuthorName string `json:"author_name"`
	CreatedAt  string `json:"created_at"`
}

var taskCmd = &cobra.Command{
	Use:   "task",
	Short: "Task management commands",
	Long:  "Commands for managing tasks",
}

var createTaskCmd = &cobra.Command{
	Use:   "create <title>",
	Short: "Create a new task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title := args[0]
		description, _ := cmd.Flags().GetString("description")
		projectID, _ := cmd.Flags().GetString("project")
		priority, _ := cmd.Flags().GetString("priority")

		c := client.New()
		req := map[string]interface{}{
			"title": title,
		}
		if description != "" {
			req["description"] = description
		}
		if projectID != "" {
			req["project_id"] = projectID
		}
		if priority != "" {
			req["priority"] = priority
		}

		resp, err := c.Post("/agent/tasks", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to create task: %s", string(body))
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Task created: %s (ID: %s)\n", task.Title, task.ID)
		return nil
	},
}

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tasks",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetString("priority")
		projectID, _ := cmd.Flags().GetString("project")
		agentID, _ := cmd.Flags().GetString("agent")
		search, _ := cmd.Flags().GetString("search")

		path := "/agent/tasks"
		params := []string{}
		if status != "" {
			params = append(params, "status="+status)
		}
		if priority != "" {
			params = append(params, "priority="+priority)
		}
		if projectID != "" {
			params = append(params, "project_id="+projectID)
		}
		if agentID != "" {
			params = append(params, "agent_id="+agentID)
		}
		if search != "" {
			params = append(params, "search="+search)
		}

		if len(params) > 0 {
			path += "?" + strings.Join(params, "&")
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
		fmt.Fprintln(w, "ID\tTITLE\tSTATUS\tPRIORITY\tPROJECT")
		for _, t := range tasks {
			projectID := t.ProjectID
			if projectID == "" {
				projectID = "-"
			}
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", t.ID, t.Title, t.Status, t.Priority, projectID)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d task(s)\n", len(tasks))
		return nil
	},
}

var getTaskCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get task details",
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
		fmt.Printf("Project ID:  %s\n", task.ProjectID)
		if task.AgentID != nil && *task.AgentID != "" {
			fmt.Printf("Agent ID:    %s\n", *task.AgentID)
		}
		fmt.Printf("Created:     %s\n", task.CreatedAt)
		fmt.Printf("Updated:     %s\n", task.UpdatedAt)
		return nil
	},
}

var updateTaskCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		title, _ := cmd.Flags().GetString("title")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		priority, _ := cmd.Flags().GetString("priority")

		req := make(map[string]interface{})
		if title != "" {
			req["title"] = title
		}
		if description != "" {
			req["description"] = description
		}
		if status != "" {
			req["status"] = status
		}
		if priority != "" {
			req["priority"] = priority
		}

		if len(req) == 0 {
			return fmt.Errorf("at least one field must be specified")
		}

		c := client.New()
		resp, err := c.Patch("/agent/tasks/"+id, req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update task: %s", string(body))
		}

		var task Task
		if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Task updated: %s\n", task.Title)
		return nil
	},
}

var deleteTaskCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Delete("/agent/tasks/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete task: %s", string(body))
		}

		fmt.Println("✓ Task deleted")
		return nil
	},
}

var assignTaskCmd = &cobra.Command{
	Use:   "assign <task-id> <agent-id>",
	Short: "Assign an agent to a task",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]
		agentID := args[1]

		c := client.New()
		req := map[string]string{
			"agent_id": agentID,
		}

		resp, err := c.Patch("/tasks/"+taskID+"/assign", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to assign task: %s", string(body))
		}

		fmt.Println("✓ Agent assigned to task")
		return nil
	},
}

var unassignTaskCmd = &cobra.Command{
	Use:   "unassign <task-id>",
	Short: "Unassign agent from a task",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		taskID := args[0]

		c := client.New()
		resp, err := c.Patch("/tasks/"+taskID+"/unassign", nil, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to unassign task: %s", string(body))
		}

		fmt.Println("✓ Agent unassigned from task")
		return nil
	},
}

var listCommentsCmd = &cobra.Command{
	Use:   "comments <task-id>",
	Short: "List task comments",
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
			fmt.Printf("[%s] %s\n", comment.CreatedAt, comment.Content)
		}

		fmt.Printf("\nTotal: %d comment(s)\n", len(comments))
		return nil
	},
}

var addCommentCmd = &cobra.Command{
	Use:   "comment <task-id> <content>",
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

		fmt.Println("✓ Comment added")
		return nil
	},
}

var completeTaskCmd = &cobra.Command{
	Use:     "complete <task-id>",
	Short:   "Mark a task as completed",
	Args:    cobra.ExactArgs(1),
	Aliases: []string{"done"},
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		comment, _ := cmd.Flags().GetString("comment")

		c := client.New()
		req := map[string]string{
			"status": "completed",
		}

		resp, err := c.Patch("/agent/tasks/"+id+"/status", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to complete task: %s", string(body))
		}

		// Add comment if provided
		if comment != "" {
			commentReq := map[string]string{"content": comment}
			c.Post("/agent/tasks/"+id+"/comments", commentReq, true)
		}

		fmt.Printf("✓ Task %s marked as completed\n", id)
		return nil
	},
}

var startTaskCmd = &cobra.Command{
	Use:   "start <task-id>",
	Short: "Mark a task as in progress",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		comment, _ := cmd.Flags().GetString("comment")

		c := client.New()
		req := map[string]string{
			"status": "in_progress",
		}

		resp, err := c.Patch("/agent/tasks/"+id+"/status", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to start task: %s", string(body))
		}

		// Add comment if provided
		if comment != "" {
			commentReq := map[string]string{"content": comment}
			c.Post("/agent/tasks/"+id+"/comments", commentReq, true)
		}

		fmt.Printf("✓ Task %s marked as in progress\n", id)
		return nil
	},
}

var blockTaskCmd = &cobra.Command{
	Use:   "block <task-id>",
	Short: "Mark a task as blocked",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		reason, _ := cmd.Flags().GetString("reason")

		c := client.New()
		req := map[string]string{
			"status": "blocked",
		}

		resp, err := c.Patch("/agent/tasks/"+id+"/status", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to block task: %s", string(body))
		}

		// Add blocker reason as comment
		if reason != "" {
			commentReq := map[string]string{"content": "Blocked: " + reason}
			c.Post("/agent/tasks/"+id+"/comments", commentReq, true)
		}

		fmt.Printf("✓ Task %s marked as blocked\n", id)
		return nil
	},
}

var statsTaskCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show task statistics and analytics",
	Long:  "Display task statistics including overview, agent metrics, and timeline data",
	RunE: func(cmd *cobra.Command, args []string) error {
		view, _ := cmd.Flags().GetString("view")
		days, _ := cmd.Flags().GetInt("days")

		c := client.New()
		
		var endpoint string
		switch view {
		case "overview", "":
			endpoint = "/analytics/tasks/overview"
		case "agents":
			endpoint = "/analytics/tasks/agents"
		case "timeline":
			endpoint = fmt.Sprintf("/analytics/tasks/timeline?days=%d", days)
		default:
			return fmt.Errorf("invalid view: %s (use: overview, agents, or timeline)", view)
		}

		resp, err := c.Get(endpoint, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get stats: %s", string(body))
		}

		var stats map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		// Format output based on view type
		switch view {
		case "overview", "":
			printOverviewStats(stats)
		case "agents":
			printAgentStats(stats)
		case "timeline":
			printTimelineStats(stats)
		}

		return nil
	},
}

func printOverviewStats(stats map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	fmt.Fprintln(w, "TASK OVERVIEW")
	fmt.Fprintln(w, "=============")
	
	if total, ok := stats["total_tasks"].(float64); ok {
		fmt.Fprintf(w, "Total Tasks:\t%d\n", int(total))
	}
	if completed, ok := stats["completed_tasks"].(float64); ok {
		fmt.Fprintf(w, "Completed:\t%d\n", int(completed))
	}
	if rate, ok := stats["completion_rate"].(float64); ok {
		fmt.Fprintf(w, "Completion Rate:\t%.1f%%\n", rate)
	}
	if recent, ok := stats["tasks_last_7_days"].(float64); ok {
		fmt.Fprintf(w, "Last 7 Days:\t%d\n", int(recent))
	}
	if avgTime, ok := stats["avg_completion_time_hours"].(float64); ok {
		fmt.Fprintf(w, "Avg Completion Time:\t%.1f hours\n", avgTime)
	}
	
	fmt.Fprintln(w)
	
	// Status breakdown
	if byStatus, ok := stats["by_status"].([]interface{}); ok && len(byStatus) > 0 {
		fmt.Fprintln(w, "\nBY STATUS")
		fmt.Fprintln(w, "---------")
		for _, s := range byStatus {
			if status, ok := s.(map[string]interface{}); ok {
				statusStr, _ := status["status"].(string)
				count, _ := status["count"].(float64)
				fmt.Fprintf(w, "%s:\t%d\n", strings.Title(statusStr), int(count))
			}
		}
	}
	
	// Priority breakdown
	if byPriority, ok := stats["by_priority"].([]interface{}); ok && len(byPriority) > 0 {
		fmt.Fprintln(w, "\nBY PRIORITY")
		fmt.Fprintln(w, "-----------")
		for _, p := range byPriority {
			if priority, ok := p.(map[string]interface{}); ok {
				priorityStr, _ := priority["priority"].(string)
				count, _ := priority["count"].(float64)
				fmt.Fprintf(w, "%s:\t%d\n", strings.Title(priorityStr), int(count))
			}
		}
	}
	
	w.Flush()
}

func printAgentStats(stats map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	fmt.Fprintln(w, "AGENT METRICS")
	fmt.Fprintln(w, "=============")
	
	if totalAgents, ok := stats["total_agents"].(float64); ok {
		fmt.Fprintf(w, "Total Agents:\t%d\n", int(totalAgents))
	}
	if activeAgents, ok := stats["active_agents"].(float64); ok {
		fmt.Fprintf(w, "Active Agents:\t%d\n", int(activeAgents))
	}
	
	// Agent breakdown
	if agents, ok := stats["agents"].([]interface{}); ok && len(agents) > 0 {
		fmt.Fprintln(w, "\nPER AGENT")
		fmt.Fprintln(w, "---------")
		for _, a := range agents {
			if agent, ok := a.(map[string]interface{}); ok {
				name, _ := agent["agent_name"].(string)
				total, _ := agent["total_tasks"].(float64)
				completed, _ := agent["completed_tasks"].(float64)
				inProgress, _ := agent["in_progress_tasks"].(float64)
				blocked, _ := agent["blocked_tasks"].(float64)
				rate, _ := agent["completion_rate"].(float64)
				avgTime, _ := agent["avg_completion_time_hours"].(float64)
				
				fmt.Fprintf(w, "\n%s:\n", name)
				fmt.Fprintf(w, "  Total:\t%d\n", int(total))
				fmt.Fprintf(w, "  Completed:\t%d\n", int(completed))
				fmt.Fprintf(w, "  In Progress:\t%d\n", int(inProgress))
				fmt.Fprintf(w, "  Blocked:\t%d\n", int(blocked))
				fmt.Fprintf(w, "  Completion Rate:\t%.1f%%\n", rate)
				fmt.Fprintf(w, "  Avg Time:\t%.1f hours\n", avgTime)
			}
		}
	}
	
	w.Flush()
}

func printTimelineStats(stats map[string]interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	
	fmt.Fprintln(w, "TIMELINE METRICS")
	fmt.Fprintln(w, "================")
	
	if avgAge, ok := stats["avg_task_age_hours"].(float64); ok {
		fmt.Fprintf(w, "Average Task Age:\t%.1f hours\n", avgAge)
	}
	
	// Created daily
	if createdDaily, ok := stats["created_daily"].([]interface{}); ok && len(createdDaily) > 0 {
		fmt.Fprintln(w, "\nTASKS CREATED (last N days)")
		fmt.Fprintln(w, "---------------------------")
		for _, d := range createdDaily {
			if day, ok := d.(map[string]interface{}); ok {
				date, _ := day["date"].(string)
				count, _ := day["count"].(float64)
				fmt.Fprintf(w, "%s:\t%d\n", date, int(count))
			}
		}
	}
	
	// Completed daily
	if completedDaily, ok := stats["completed_daily"].([]interface{}); ok && len(completedDaily) > 0 {
		fmt.Fprintln(w, "\nTASKS COMPLETED (last N days)")
		fmt.Fprintln(w, "-----------------------------")
		for _, d := range completedDaily {
			if day, ok := d.(map[string]interface{}); ok {
				date, _ := day["date"].(string)
				count, _ := day["count"].(float64)
				fmt.Fprintf(w, "%s:\t%d\n", date, int(count))
			}
		}
	}
	
	w.Flush()
}

func init() {
	rootCmd.AddCommand(taskCmd)

	// Create task
	createTaskCmd.Flags().StringP("description", "d", "", "Task description")
	createTaskCmd.Flags().StringP("project", "p", "", "Project ID")
	createTaskCmd.Flags().StringP("priority", "P", "", "Task priority (low, medium, high)")

	// List tasks
	listTasksCmd.Flags().StringP("status", "s", "", "Filter by status")
	listTasksCmd.Flags().StringP("priority", "P", "", "Filter by priority")
	listTasksCmd.Flags().StringP("project", "p", "", "Filter by project ID")
	listTasksCmd.Flags().StringP("agent", "a", "", "Filter by agent ID")
	listTasksCmd.Flags().StringP("search", "q", "", "Search by title or description")

	// Update task
	updateTaskCmd.Flags().StringP("title", "t", "", "Task title")
	updateTaskCmd.Flags().StringP("description", "d", "", "Task description")
	updateTaskCmd.Flags().StringP("status", "s", "", "Task status")
	updateTaskCmd.Flags().StringP("priority", "P", "", "Task priority")

	taskCmd.AddCommand(createTaskCmd)
	taskCmd.AddCommand(listTasksCmd)
	taskCmd.AddCommand(getTaskCmd)
	taskCmd.AddCommand(updateTaskCmd)
	taskCmd.AddCommand(deleteTaskCmd)
	taskCmd.AddCommand(assignTaskCmd)
	taskCmd.AddCommand(unassignTaskCmd)
	taskCmd.AddCommand(listCommentsCmd)
	taskCmd.AddCommand(addCommentCmd)
	taskCmd.AddCommand(completeTaskCmd)
	taskCmd.AddCommand(startTaskCmd)
	taskCmd.AddCommand(blockTaskCmd)
	taskCmd.AddCommand(statsTaskCmd)

	// Complete task flags
	completeTaskCmd.Flags().StringP("comment", "c", "", "Optional completion comment")

	// Start task flags
	startTaskCmd.Flags().StringP("comment", "c", "", "Optional start comment")

	// Block task flags
	blockTaskCmd.Flags().StringP("reason", "r", "", "Reason for blocking")

	// Stats task flags
	statsTaskCmd.Flags().StringP("view", "v", "overview", "View type (overview, agents, timeline)")
	statsTaskCmd.Flags().IntP("days", "d", 30, "Number of days for timeline view")
}
