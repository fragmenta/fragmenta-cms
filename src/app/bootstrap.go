package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/fragmenta/query"
)

// TODO: This should probably go into a bootstrap package within fragmenta?

const (
	fragmentaVersion = "1.2"

	permissions                 = 0744
	createDatabaseMigrationName = "Create-Database"
	createTablesMigrationName   = "Create-Tables"
)

var (
	// ConfigDevelopment holds the development config from fragmenta.json
	ConfigDevelopment map[string]string

	// ConfigProduction holds development config from fragmenta.json
	ConfigProduction map[string]string

	// ConfigTest holds the app test config from fragmenta.json
	ConfigTest map[string]string
)

// Bootstrap generates missing config files, sql migrations, and runs the first migrations
// For this we need to know what to call the app, but we default to fragmenta-cms for now
// we could use our current folder name?
func Bootstrap() error {
	// We assume we're being run from root of project path
	projectPath, err := os.Getwd()
	if err != nil {
		return err
	}

	fmt.Printf("\nBootstrapping server...\n")

	err = generateConfig(projectPath)
	if err != nil {
		return err
	}

	err = generateCreateSQL(projectPath)
	if err != nil {
		return err
	}

	// Run the migrations without the fragmenta tool being present
	err = runMigrations(projectPath)
	if err != nil {
		return err
	}

	return nil
}

// RequiresBootStrap returns true if the app requires bootstrapping
func RequiresBootStrap() bool {
	if !fileExists(configPath()) {
		return true
	}
	return false
}

func configPath() string {
	return "secrets/fragmenta.json"
}

func projectPathRelative(projectPath string) string {
	goSrc := os.Getenv("GOPATH") + "/src/"
	return strings.Replace(projectPath, goSrc, "", 1)
}

func generateConfig(projectPath string) error {
	configPath := configPath()
	prefix := path.Base(projectPath)
	log.Printf("Generating new config at %s", configPath)

	ConfigProduction = map[string]string{}
	ConfigDevelopment = map[string]string{}
	ConfigTest = map[string]string{
		"port":            "3000",
		"log":             "log/test.log",
		"db_adapter":      "postgres",
		"db":              prefix + "_test",
		"db_user":         prefix + "_server",
		"db_pass":         randomKey(8),
		"assets_compiled": "no",
		"path":            projectPathRelative(projectPath),
		"hmac_key":        randomKey(32),
		"secret_key":      randomKey(32),
	}

	// Should we ask for db prefix when setting up?
	// hmm, in fact can we do this setup here at all!!
	for k, v := range ConfigTest {
		ConfigDevelopment[k] = v
		ConfigProduction[k] = v
	}
	ConfigDevelopment["db"] = prefix + "_development"
	ConfigDevelopment["log"] = "log/development.log"
	ConfigDevelopment["hmac_key"] = randomKey(32)
	ConfigDevelopment["secret_key"] = randomKey(32)

	ConfigProduction["db"] = prefix + "_production"
	ConfigProduction["log"] = "log/production.log"
	ConfigProduction["port"] = "80"
	ConfigProduction["assets_compiled"] = "yes"
	ConfigProduction["hmac_key"] = randomKey(32)
	ConfigProduction["secret_key"] = randomKey(32)

	configs := map[string]map[string]string{
		"production":  ConfigProduction,
		"development": ConfigDevelopment,
		"test":        ConfigTest,
	}

	configJSON, err := json.MarshalIndent(configs, "", "\t")
	if err != nil {
		log.Printf("Error parsing config %s %v", configPath, err)
		return err
	}

	// Write the config json file
	err = ioutil.WriteFile(configPath, configJSON, permissions)
	if err != nil {
		log.Printf("Error writing config %s %v", configPath, err)
		return err
	}

	return nil
}

// generateCreateSQL generates an SQL migration file to create the database user and database referred to in config
func generateCreateSQL(projectPath string) error {

	// Set up a Create-Database migration, which comes first
	name := path.Base(projectPath)
	d := ConfigDevelopment["db"]
	u := ConfigDevelopment["db_user"]
	p := ConfigDevelopment["db_pass"]
	sql := fmt.Sprintf("/* Setup database for %s */\nCREATE USER \"%s\" WITH PASSWORD '%s';\nCREATE DATABASE \"%s\" WITH OWNER \"%s\";", name, u, p, d, u)

	// Generate a migration to create db with today's date
	file := migrationPath(projectPath, createDatabaseMigrationName)
	err := ioutil.WriteFile(file, []byte(sql), 0744)
	if err != nil {
		return err
	}

	// If we have a Create-Tables file, copy it out to a new migration with today's date
	createTablesPath := path.Join(projectPath, "db", "migrate", createTablesMigrationName+".sql.tmpl")
	if fileExists(createTablesPath) {
		sql, err := ioutil.ReadFile(createTablesPath)
		if err != nil {
			return err
		}

		// Now vivify the template, for now we just replace one key
		sqlString := strings.Replace(string(sql), "[[.fragmenta_db_user]]", u, -1)

		file = migrationPath(projectPath, createTablesMigrationName)
		err = ioutil.WriteFile(file, []byte(sqlString), 0744)
		if err != nil {
			return err
		}
		// Remove the old file
		os.Remove(createTablesPath)

	} else {
		fmt.Printf("NO TABLES %s", createTablesPath)
	}

	return nil
}

// runMigrations at projectPath
func runMigrations(projectPath string) error {
	var migrations []string
	var migrationCount int

	config := ConfigDevelopment

	// Get a list of migration files
	files, err := filepath.Glob("./db/migrate/*.sql")
	if err != nil {
		return err
	}

	// Sort the list alphabetically
	sort.Strings(files)

	for _, file := range files {
		filename := path.Base(file)

		log.Printf("Running migration %s", filename)

		args := []string{"-d", config["db"], "-f", file}
		if strings.Contains(filename, createDatabaseMigrationName) {
			args = []string{"-f", file}
			log.Printf("Running database creation migration: %s", file)
		}

		// Execute this sql file against the database
		result, err := runCommand("psql", args...)
		if err != nil || strings.Contains(string(result), "ERROR") {
			if err == nil {
				err = fmt.Errorf("\n%s", string(result))
			}

			// If at any point we fail, log it and break
			log.Printf("ERROR loading sql migration:%s\n", err)
			log.Printf("All further migrations cancelled\n\n")
			return err
		}

		migrationCount++
		migrations = append(migrations, filename)
		log.Printf("Completed migration %s\n%s\n%s", filename, string(result), "-")

	}

	if migrationCount > 0 {
		writeMetadata(config, migrations)
		log.Printf("Migrations complete up to migration %v on db %s\n\n", migrations, config["db"])
	}

	return nil
}

// Oh, we need to write the full list of migrations, not just one migration version

// Update the database with a line recording what we have done
func writeMetadata(config map[string]string, migrations []string) {
	// Try opening the db (db may not exist at this stage)
	err := openDatabase(config)
	if err != nil {
		log.Printf("Database ERROR %s", err)
	}
	defer query.CloseDatabase()

	for _, m := range migrations {
		sql := "Insert into fragmenta_metadata(updated_at,fragmenta_version,migration_version,status) VALUES(NOW(),$1,$2,100);"
		result, err := query.ExecSQL(sql, fragmentaVersion, m)
		if err != nil {
			log.Printf("Database ERROR %s %s", err, result)
		}
	}

}

// Open our database
func openDatabase(config map[string]string) error {
	// Open the database
	options := map[string]string{
		"adapter":  config["db_adapter"],
		"user":     config["db_user"],
		"password": config["db_pass"],
		"db":       config["db"],
		// "debug"     : "true",
	}

	err := query.OpenDatabase(options)
	if err != nil {
		return err
	}

	log.Printf("%s\n", "-")
	log.Printf("Opened database at %s for user %s", config["db"], config["db_user"])
	return nil
}

// Generate a suitable path for a migration from the current date/time down to nanosecond
func migrationPath(path string, name string) string {
	now := time.Now()
	layout := "2006-01-02-150405"
	return fmt.Sprintf("%s/db/migrate/%s-%s.sql", path, now.Format(layout), name)
}

// Generate a random 32 byte key encoded in base64
func randomKey(l int64) string {
	k := make([]byte, l)
	if _, err := io.ReadFull(rand.Reader, k); err != nil {
		return ""
	}
	return hex.EncodeToString(k)
}

// fileExists returns true if this file exists
func fileExists(p string) bool {
	_, err := os.Stat(p)
	if err != nil && os.IsNotExist(err) {
		return false
	}

	return true
}

// runCommand runs a command with exec.Command
func runCommand(command string, args ...string) ([]byte, error) {

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, err
	}

	return output, nil
}
