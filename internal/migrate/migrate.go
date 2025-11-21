package migrate

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/e"
	"github.com/SALutHere/avito-2025-autumn-backend-internship/pkg/logger"
)

const migrationsDir = "migrations"

func Run(db *sql.DB) error {
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return e.Wrap("can not read migrations folder", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}

	sort.Strings(sqlFiles)

	for _, fileName := range sqlFiles {
		path := filepath.Join("migrations", fileName)

		logger.L().Info("Applying migration", slog.String("file", fileName))

		content, err := os.ReadFile(path)
		if err != nil {
			return e.Wrap(fmt.Sprintf("can not read %s", fileName), err)
		}

		if _, err = db.Exec(string(content)); err != nil {
			return e.Wrap(fmt.Sprintf("can not execute %s", fileName), err)
		}
	}

	logger.L().Info("All migrations applied successfully")
	return nil
}
