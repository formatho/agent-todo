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

		resp, err := c.Post("/tasks", req, true)
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

		path := "/tasks"
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
		resp, err := c.Get("/tasks/"+id, true)
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
		resp, err := c.Patch("/tasks/"+id, req, true)
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
		resp, err := c.Delete("/tasks/"+id, true)
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
		resp, err := c.Get("/tasks/"+taskID+"/comments", true)
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

		resp, err := c.Post("/tasks/"+taskID+"/comments", req, true)
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
}
