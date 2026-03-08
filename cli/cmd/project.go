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

type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Project management commands",
	Long:  "Commands for managing projects",
}

var createProjectCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		c := client.New()
		req := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			req["description"] = description
		}
		if status != "" {
			req["status"] = status
		}

		resp, err := c.Post("/projects", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to create project: %s", string(body))
		}

		var project Project
		if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Project created: %s (ID: %s)\n", project.Name, project.ID)
		return nil
	},
}

var listProjectsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		search, _ := cmd.Flags().GetString("search")

		path := "/projects"
		if status != "" || search != "" {
			path += "?"
			if status != "" {
				path += "status=" + status
				if search != "" {
					path += "&"
				}
			}
			if search != "" {
				path += "search=" + search
			}
		}

		c := client.New()
		resp, err := c.Get(path, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to list projects: %s", string(body))
		}

		var projects []Project
		if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		if len(projects) == 0 {
			fmt.Println("No projects found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDESCRIPTION")
		for _, p := range projects {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", p.ID, p.Name, p.Status, p.Description)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d project(s)\n", len(projects))
		return nil
	},
}

var getProjectCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get project details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Get("/projects/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get project: %s", string(body))
		}

		var project Project
		if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("ID:          %s\n", project.ID)
		fmt.Printf("Name:        %s\n", project.Name)
		fmt.Printf("Description: %s\n", project.Description)
		fmt.Printf("Status:      %s\n", project.Status)
		fmt.Printf("Created:     %s\n", project.CreatedAt)
		fmt.Printf("Updated:     %s\n", project.UpdatedAt)
		return nil
	},
}

var updateProjectCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")

		req := make(map[string]interface{})
		if name != "" {
			req["name"] = name
		}
		if description != "" {
			req["description"] = description
		}
		if status != "" {
			req["status"] = status
		}

		if len(req) == 0 {
			return fmt.Errorf("at least one field must be specified")
		}

		c := client.New()
		resp, err := c.Patch("/projects/"+id, req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update project: %s", string(body))
		}

		var project Project
		if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Project updated: %s\n", project.Name)
		return nil
	},
}

var deleteProjectCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Delete("/projects/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete project: %s", string(body))
		}

		fmt.Println("✓ Project deleted")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Create project
	createProjectCmd.Flags().StringP("description", "d", "", "Project description")
	createProjectCmd.Flags().StringP("status", "s", "", "Project status")

	// List projects
	listProjectsCmd.Flags().StringP("status", "s", "", "Filter by status")
	listProjectsCmd.Flags().StringP("search", "q", "", "Search by name or description")

	// Update project
	updateProjectCmd.Flags().StringP("name", "n", "", "Project name")
	updateProjectCmd.Flags().StringP("description", "d", "", "Project description")
	updateProjectCmd.Flags().StringP("status", "s", "", "Project status")

	projectCmd.AddCommand(createProjectCmd)
	projectCmd.AddCommand(listProjectsCmd)
	projectCmd.AddCommand(getProjectCmd)
	projectCmd.AddCommand(updateProjectCmd)
	projectCmd.AddCommand(deleteProjectCmd)
}
