package config

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrationsFolder runs all *.sql files in the specified folder in sorted order.
func RunMigrationsFolder(db *sql.DB, folder string) error {
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		return fmt.Errorf("cannot read migrations dir: %w", err)
	}

	var sqlFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles) // Ensure 001, 002, 003 order

	for _, fname := range sqlFiles {
		fullpath := filepath.Join(folder, fname)
		sqlBytes, err := ioutil.ReadFile(fullpath)
		if err != nil {
			return fmt.Errorf("cannot read migration file %s: %w", fname, err)
		}
		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("error running migration %s: %w", fname, err)
		}
		fmt.Printf("Applied migration: %s\n", fname)
	}
	return nil
}
