package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/formatho/agent-todo/cli/client"
	"github.com/formatho/agent-todo/cli/config"
	"github.com/spf13/cobra"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name,omitempty"`
}

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	CreatedAt   string `json:"created_at"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authentication commands",
	Long:  "Commands for user authentication (login, register, logout, whoami)",
}

var loginCmd = &cobra.Command{
	Use:   "login <email> <password>",
	Short: "Login to the platform",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		password := args[1]

		c := client.New()
		req := LoginRequest{
			Email:    email,
			Password: password,
		}

		resp, err := c.Post("/auth/login", req, false)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("login failed: %s", string(body))
		}

		var authResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		// Save token to config
		config.SetToken(authResp.Token)
		if err := config.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Printf("✓ Logged in successfully as %s (%s)\n", authResp.User.Email, authResp.User.DisplayName)
		return nil
	},
}

var registerCmd = &cobra.Command{
	Use:   "register <email> <password> [display-name]",
	Short: "Register a new account",
	Args:  cobra.RangeArgs(2, 3),
	RunE: func(cmd *cobra.Command, args []string) error {
		email := args[0]
		password := args[1]
		displayName := ""
		if len(args) > 2 {
			displayName = args[2]
		}

		c := client.New()
		req := RegisterRequest{
			Email:       email,
			Password:    password,
			DisplayName: displayName,
		}

		resp, err := c.Post("/auth/register", req, false)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("registration failed: %s", string(body))
		}

		var authResp AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		// Save token to config
		config.SetToken(authResp.Token)
		if err := config.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Printf("✓ Registered and logged in successfully as %s\n", authResp.User.Email)
		return nil
	},
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from the platform",
	RunE: func(cmd *cobra.Command, args []string) error {
		config.ClearAuth()
		if err := config.Save(); err != nil {
			return fmt.Errorf("error saving config: %w", err)
		}

		fmt.Println("✓ Logged out successfully")
		return nil
	},
}

var whoamiCmd = &cobra.Command{
	Use:   "whoami",
	Short: "Show current user information",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := config.Get()
		if cfg.Token == "" {
			return fmt.Errorf("not logged in. Use 'agent-todo auth login' first")
		}

		c := client.New()
		resp, err := c.Get("/auth/me", true)
		if err != nil {
			return fmt.Errorf("error making request: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return fmt.Errorf("error getting user info: %s", string(body))
		}

		var user User
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}

		fmt.Printf("Email: %s\n", user.Email)
		fmt.Printf("Name: %s\n", user.DisplayName)
		fmt.Printf("ID: %s\n", user.ID)
		fmt.Printf("Joined: %s\n", user.CreatedAt)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(logoutCmd)
	authCmd.AddCommand(whoamiCmd)
}
