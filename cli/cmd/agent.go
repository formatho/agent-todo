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

	agentCmd.AddCommand(createAgentCmd)
	agentCmd.AddCommand(listAgentsCmd)
	agentCmd.AddCommand(getAgentCmd)
	agentCmd.AddCommand(updateAgentCmd)
	agentCmd.AddCommand(deleteAgentCmd)
}
