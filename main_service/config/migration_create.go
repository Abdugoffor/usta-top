package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// make migrate-create name=add_code_to_regions
// # ✅ Migration created: config/migrations/002_add_code_to_regions.sql

const migrationsDir = "config/migrations"

func MigrateCreate(name string) {
	name = strings.ToLower(strings.TrimSpace(name))
	name = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(name, "_")

	if name == "" {
		fmt.Println("❌ Migration name cannot be empty")
		os.Exit(1)
	}

	next := nextMigrationNumber()
	fullPath := filepath.Join(migrationsDir, fmt.Sprintf("%03d_%s.sql", next, name))

	if err := os.WriteFile(fullPath, []byte(""), 0644); err != nil {
		fmt.Printf("❌ Failed to create file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Migration created: %s\n", fullPath)
}

func nextMigrationNumber() int {
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Printf("❌ Cannot read migrations dir: %v\n", err)
		os.Exit(1)
	}

	re := regexp.MustCompile(`^(\d+)_`)
	nums := []int{}

	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		if m := re.FindStringSubmatch(e.Name()); len(m) > 1 {
			if n, err := strconv.Atoi(m[1]); err == nil {
				nums = append(nums, n)
			}
		}
	}

	if len(nums) > 0 {
		sort.Ints(nums)
		return nums[len(nums)-1] + 1
	}

	return 1
}
