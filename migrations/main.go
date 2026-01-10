package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

var (
	flags        = flag.NewFlagSet("goose", flag.ExitOnError)
	dirname      = flags.String("dir", "./migrations", "directory with migration files")
	tableName    = flags.String("table", "goose_db_version", "migration table name")
	verbose      = flags.Bool("v", false, "enable verbose mode")
	help         = flags.Bool("h", false, "print help")
	version      = flags.Bool("version", false, "print version")
	sequential   = flags.Bool("s", false, "use sequential numbering for new migrations")
	allowMissing = flags.Bool("allow-missing", false, "applies missing (out-of-order) migrations")
)

func main() {
	// Try to load .env file if it exists (silent fail if not found)
	_ = godotenv.Load()

	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if *version {
		fmt.Println("goose version:", goose.Version)
		return
	}

	if *help {
		flags.Usage()
		return
	}

	args := flags.Args()
	if len(args) == 0 {
		flags.Usage()
		os.Exit(1)
	}

	// Get database connection string from environment or use defaults
	dbString := getDBString()

	if dbString == "" {
		log.Fatal("database connection string is required. Set DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME environment variables")
	}

	// Normalize directory path
	dir := *dirname

	// Open database connection
	db, err := sql.Open("pgx", dbString)
	if err != nil {
		log.Fatalf("goose: failed to open database: %v", err)
	}
	defer db.Close()

	// Set options
	goose.SetTableName(*tableName)
	if *verbose {
		goose.SetVerbose(true)
	}
	if *sequential {
		goose.SetSequential(true)
	}

	// Get command
	command := args[0]
	commandArgs := args[1:]

	// Special handling for bootstrap command (custom command)
	if command == "bootstrap" {
		if err := bootstrap(db); err != nil {
			log.Fatalf("bootstrap failed: %v", err)
		}
		return
	}

	// Execute command
	if err := goose.Run(command, db, dir, commandArgs...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}

// getDBString constructs the database connection string from environment variables
func getDBString() string {
	host := os.Getenv("DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("DB_USER")
	if user == "" {
		user = "kg_user"
	}

	password := os.Getenv("DB_PASSWORD")
	if password == "" {
		return "" // Require password
	}

	dbname := os.Getenv("DB_NAME")
	if dbname == "" {
		dbname = "knowledge_garden"
	}

	sslmode := os.Getenv("DB_SSL_MODE")
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)
}

// bootstrap marks all existing migrations as applied without running them
// Use this for databases that were set up before the migration system
func bootstrap(db *sql.DB) error {
	// Use the default migrations directory
	migrationsDir := "./migrations"

	// Create goose tracking table if it doesn't exist
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS goose_db_version (
			id SERIAL PRIMARY KEY,
			version_id BIGINT NOT NULL,
			is_applied BOOLEAN NOT NULL,
			tstamp TIMESTAMP NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create goose_db_version table: %w", err)
	}

	fmt.Println("Bootstrapping migration tracking for existing database...")
	fmt.Println("This will mark all current migrations as applied (without running SQL).")

	// Read migration files and extract version numbers
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Process each migration file
	for _, entry := range entries {
		// Skip non-SQL files and main.go
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") || entry.Name() == "main.go" {
			continue
		}

		// Extract version number from filename (format: YYYYMMDDHHMMSS_description.sql)
		filename := entry.Name()
		parts := strings.SplitN(filename, "_", 2)
		if len(parts) < 2 {
			continue
		}

		// Parse version ID (full timestamp is the version)
		version, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			log.Printf("Warning: skipping file with invalid version: %s", filename)
			continue
		}

		// Check if already applied
		var count int
		err = db.QueryRow("SELECT COUNT(*) FROM goose_db_version WHERE version_id = $1", version).Scan(&count)
		if err != nil {
			return fmt.Errorf("failed to check migration %d: %w", version, err)
		}

		if count > 0 {
			fmt.Printf("  Migration %d already tracked - skipping\n", version)
			continue
		}

		// Insert as applied
		_, err = db.Exec("INSERT INTO goose_db_version (version_id, is_applied, tstamp) VALUES ($1, true, NOW())", version)
		if err != nil {
			return fmt.Errorf("failed to insert migration %d: %w", version, err)
		}
		fmt.Printf("  ✓ Marked migration %d as applied: %s\n", version, filename)
	}

	fmt.Println("\n✓ Bootstrap complete!")
	fmt.Println("  All migrations are now tracked. Future migrations will run normally.")
	fmt.Println("  Note: Migrations are idempotent and can be re-run safely.")
	return nil
}

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  go run ./migrations [options] <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up                   Migrate all pending migrations")
	fmt.Println("  up-by-one            Migrate only the next pending migration")
	fmt.Println("  up-to VERSION        Migrate all pending migrations up to VERSION")
	fmt.Println("  down                 Rollback the most recent migration")
	fmt.Println("  down-to VERSION      Rollback all migrations down to VERSION")
	fmt.Println("  redo                 Rollback the most recent migration then re-apply it")
	fmt.Println("  reset                Rollback all migrations and re-apply them all")
	fmt.Println("  status               Print the status of all migrations")
	fmt.Println("  version              Print the current migration version")
	fmt.Println("  create NAME [type]   Creates new migration file with NAME and optional TYPE (sql by default)")
	fmt.Println("  bootstrap            Mark all migrations as applied (for existing databases)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -dir string         Directory with migration files (default \"./migrations\")")
	fmt.Println("  -table string       Migration table name (default \"goose_db_version\")")
	fmt.Println("  -v                  Enable verbose mode")
	fmt.Println("  -h                  Print help")
	fmt.Println("  -version            Print version")
	fmt.Println("  -s                  Use sequential numbering for new migrations")
	fmt.Println("  -allow-missing      Applies missing (out-of-order) migrations")
	fmt.Println()
	fmt.Println("Environment Variables:")
	fmt.Println("  DB_HOST             Database host (default: localhost)")
	fmt.Println("  DB_PORT             Database port (default: 5432)")
	fmt.Println("  DB_USER             Database user (default: kg_user)")
	fmt.Println("  DB_PASSWORD         Database password (required)")
	fmt.Println("  DB_NAME             Database name (default: knowledge_garden)")
	fmt.Println("  DB_SSL_MODE         SSL mode (default: disable)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run ./migrations status")
	fmt.Println("  go run ./migrations up")
	fmt.Println("  go run ./migrations create add_users_table sql")
	fmt.Println("  go run ./migrations down")
	fmt.Println("  go run ./migrations bootstrap    # For existing databases")
}
