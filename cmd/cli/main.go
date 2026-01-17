package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/term"

	"github.com/momokii/go-cli-notes/cmd/cli/client"
	"github.com/spf13/cobra"
)

// Version is set at build time using ldflags
var Version = "dev"

// readPassword reads password from terminal without echoing input
func readPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add newline after password input
	return string(password), err
}

var (
	config = &Config{}
	apiClient *client.APIClient
	authState = &client.AuthState{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "kg-cli",
	Short:   "Personal Knowledge Garden CLI",
	Version: Version,
	Long: `A terminal-based note-taking application with wiki-style links,
full-text search, and knowledge graph visualization.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Load configuration
		cfg, err := LoadConfig()
		if err != nil {
			return fmt.Errorf("load config: %w", err)
		}
		config = cfg

		// Initialize API client
		apiClient = client.NewAPIClient(
			cfg.API.BaseURL,
			time.Duration(cfg.API.Timeout)*time.Second,
		)

		// Load authentication state
		state, err := client.LoadAuthState()
		if err != nil {
			return fmt.Errorf("load auth state: %w", err)
		}
		authState = state

		// Apply auth state to client if available
		if state.IsAuthenticated() {
			state.ApplyToClient(apiClient)
		}

		return nil
	},
}

// loginCmd handles user login
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Login to your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		var email string

		fmt.Print("Email: ")
		fmt.Scanln(&email)

		password, err := readPassword("Password: ")
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}

		if email == "" || password == "" {
			return fmt.Errorf("email and password are required")
		}

		// Attempt login
		authResp, err := apiClient.Login(email, password)
		if err != nil {
			return fmt.Errorf("login failed: %w", err)
		}

		// Save auth state
		authState = &client.AuthState{
			AccessToken:  authResp.AccessToken,
			RefreshToken: authResp.RefreshToken,
			Email:        email,
		}

		if err := client.SaveAuthState(authState); err != nil {
			return fmt.Errorf("save auth state: %w", err)
		}

		fmt.Println("Login successful!")
		return nil
	},
}

// registerCmd handles user registration
var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new account",
	RunE: func(cmd *cobra.Command, args []string) error {
		var username, email string

		fmt.Print("Username: ")
		fmt.Scanln(&username)

		fmt.Print("Email: ")
		fmt.Scanln(&email)

		// Password with confirmation loop
		var password string
		for {
			pw, err := readPassword("Password: ")
			if err != nil {
				return fmt.Errorf("read password: %w", err)
			}

			confirm, err := readPassword("Confirm Password: ")
			if err != nil {
				return fmt.Errorf("read password: %w", err)
			}

			if pw == confirm {
				password = pw
				break
			}

			fmt.Println("Passwords do not match. Please try again.")
		}

		if username == "" || email == "" || password == "" {
			return fmt.Errorf("username, email and password are required")
		}

		// Attempt registration
		if err := apiClient.Register(username, email, password); err != nil {
			return fmt.Errorf("registration failed: %w", err)
		}

		fmt.Println("Registration successful! Please login with `kg-cli login`")
		return nil
	},
}

// logoutCmd handles user logout
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Logout from your account",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !authState.IsAuthenticated() {
			fmt.Println("Not logged in")
			return nil
		}

		// Call API logout
		if err := apiClient.Logout(); err != nil {
			fmt.Printf("Warning: API logout failed: %v\n", err)
		}

		// Clear local auth state
		if err := client.ClearAuthState(); err != nil {
			return fmt.Errorf("clear auth state: %w", err)
		}

		fmt.Println("Logged out successfully")
		return nil
	},
}

// statusCmd shows authentication status
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show authentication and connection status",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Knowledge Garden CLI Status")
		fmt.Println("==========================")

		// Config info
		fmt.Printf("API URL: %s\n", config.API.BaseURL)

		// Auth status - validate with server
		if authState.IsAuthenticated() {
			// Check if token is actually valid by calling the server
			if apiClient.ValidateToken() {
				fmt.Println("Status: Authenticated")
				if authState.Email != "" {
					fmt.Printf("Email: %s\n", authState.Email)
				}
			} else {
				// Token exists but is invalid/expired
				fmt.Println("Status: Not authenticated (token expired)")
				fmt.Println("\nYour session has expired. Please run 'kg-cli login' to authenticate.")
			}
		} else {
			fmt.Println("Status: Not authenticated")
			fmt.Println("\nUse 'kg-cli login' to authenticate")
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add subcommands
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(registerCmd)
	rootCmd.AddCommand(logoutCmd)
	rootCmd.AddCommand(statusCmd)
}

func main() {
	if err := Execute(); err != nil {
		os.Exit(1)
	}
}
