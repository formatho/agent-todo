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

type Organisation struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type OrganisationMember struct {
	ID       string `json:"id"`
	UserID   string `json:"user_id"`
	OrgID    string `json:"organisation_id"`
	Role     string `json:"role"`
	JoinedAt string `json:"joined_at"`
	User     struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
}

var organisationCmd = &cobra.Command{
	Use:     "organisation",
	Short:   "Organisation management commands",
	Long:    "Commands for managing organisations",
	Aliases: []string{"org"},
}

var createOrgCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create a new organisation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		description, _ := cmd.Flags().GetString("description")

		c := client.New()
		req := map[string]interface{}{
			"name": name,
		}
		if description != "" {
			req["description"] = description
		}

		resp, err := c.Post("/organisations", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to create organisation: %s", string(body))
		}

		var org Organisation
		if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Organisation created: %s (ID: %s)\n", org.Name, org.ID)
		return nil
	},
}

var listOrgsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all organisations you belong to",
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		search, _ := cmd.Flags().GetString("search")

		path := "/organisations"
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
			return fmt.Errorf("failed to list organisations: %s", string(body))
		}

		var orgs []Organisation
		if err := json.NewDecoder(resp.Body).Decode(&orgs); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		if len(orgs) == 0 {
			fmt.Println("No organisations found")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
		fmt.Fprintln(w, "ID\tNAME\tSTATUS\tDESCRIPTION")
		for _, o := range orgs {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", o.ID, o.Name, o.Status, o.Description)
		}
		w.Flush()

		fmt.Printf("\nTotal: %d organisation(s)\n", len(orgs))
		return nil
	},
}

var getOrgCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "Get organisation details",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Get("/organisations/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to get organisation: %s", string(body))
		}

		var org Organisation
		if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("ID:          %s\n", org.ID)
		fmt.Printf("Name:        %s\n", org.Name)
		fmt.Printf("Description: %s\n", org.Description)
		fmt.Printf("Status:      %s\n", org.Status)
		fmt.Printf("Created:     %s\n", org.CreatedAt)
		fmt.Printf("Updated:     %s\n", org.UpdatedAt)
		return nil
	},
}

var updateOrgCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update an organisation",
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
		resp, err := c.Patch("/organisations/"+id, req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update organisation: %s", string(body))
		}

		var org Organisation
		if err := json.NewDecoder(resp.Body).Decode(&org); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Organisation updated: %s\n", org.Name)
		return nil
	},
}

var deleteOrgCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete an organisation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		c := client.New()
		resp, err := c.Delete("/organisations/"+id, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to delete organisation: %s", string(body))
		}

		fmt.Println("✓ Organisation deleted")
		return nil
	},
}

var addMemberCmd = &cobra.Command{
	Use:   "add-member <org-id>",
	Short: "Add a member to organisation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID := args[0]
		email, _ := cmd.Flags().GetString("email")
		role, _ := cmd.Flags().GetString("role")

		if email == "" {
			return fmt.Errorf("email is required")
		}
		if role == "" {
			return fmt.Errorf("role is required")
		}

		req := map[string]interface{}{
			"email": email,
			"role":  role,
		}

		c := client.New()
		resp, err := c.Post("/organisations/"+orgID+"/members", req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to add member: %s", string(body))
		}

		var member OrganisationMember
		if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Member added: %s (%s)\n", member.User.Email, member.Role)
		return nil
	},
}

var updateMemberRoleCmd = &cobra.Command{
	Use:   "update-member <org-id> <member-id>",
	Short: "Update a member's role",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID := args[0]
		memberID := args[1]
		role, _ := cmd.Flags().GetString("role")

		if role == "" {
			return fmt.Errorf("role is required")
		}

		req := map[string]interface{}{
			"role": role,
		}

		c := client.New()
		resp, err := c.Patch("/organisations/"+orgID+"/members/"+memberID, req, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to update member role: %s", string(body))
		}

		var member OrganisationMember
		if err := json.NewDecoder(resp.Body).Decode(&member); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("✓ Member role updated: %s -> %s\n", member.User.Email, member.Role)
		return nil
	},
}

var removeMemberCmd = &cobra.Command{
	Use:   "remove-member <org-id> <member-id>",
	Short: "Remove a member from organisation",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID := args[0]
		memberID := args[1]

		c := client.New()
		resp, err := c.Delete("/organisations/"+orgID+"/members/"+memberID, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to remove member: %s", string(body))
		}

		fmt.Println("✓ Member removed")
		return nil
	},
}

var leaveOrgCmd = &cobra.Command{
	Use:   "leave <org-id>",
	Short: "Leave an organisation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		orgID := args[0]

		c := client.New()
		resp, err := c.Post("/organisations/"+orgID+"/leave", nil, true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("failed to leave organisation: %s", string(body))
		}

		fmt.Println("✓ Left organisation")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(organisationCmd)

	// Create organisation
	createOrgCmd.Flags().StringP("description", "d", "", "Organisation description")

	// List organisations
	listOrgsCmd.Flags().StringP("status", "s", "", "Filter by status")
	listOrgsCmd.Flags().StringP("search", "q", "", "Search by name or description")

	// Update organisation
	updateOrgCmd.Flags().StringP("name", "n", "", "Organisation name")
	updateOrgCmd.Flags().StringP("description", "d", "", "Organisation description")
	updateOrgCmd.Flags().StringP("status", "s", "", "Organisation status")

	// Add member
	addMemberCmd.Flags().StringP("email", "e", "", "User email")
	addMemberCmd.Flags().StringP("role", "r", "", "Member role (owner, admin, member)")
	addMemberCmd.MarkFlagRequired("email")
	addMemberCmd.MarkFlagRequired("role")

	// Update member role
	updateMemberRoleCmd.Flags().StringP("role", "r", "", "New role (owner, admin, member)")
	updateMemberRoleCmd.MarkFlagRequired("role")

	organisationCmd.AddCommand(createOrgCmd)
	organisationCmd.AddCommand(listOrgsCmd)
	organisationCmd.AddCommand(getOrgCmd)
	organisationCmd.AddCommand(updateOrgCmd)
	organisationCmd.AddCommand(deleteOrgCmd)
	organisationCmd.AddCommand(addMemberCmd)
	organisationCmd.AddCommand(updateMemberRoleCmd)
	organisationCmd.AddCommand(removeMemberCmd)
	organisationCmd.AddCommand(leaveOrgCmd)
}
