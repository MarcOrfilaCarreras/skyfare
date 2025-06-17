package db

import (
    "database/sql"
    "fmt"

    _ "database/sql"
    "github.com/MarcOrfilaCarreras/skyfare/internal/db/migrations"
    "github.com/MarcOrfilaCarreras/skyfare/internal/logging"
)

type DatabaseMigrator struct {
    DB *sql.DB
}

func NewDatabaseMigrator(dbPath string) (*DatabaseMigrator, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, err
    }
    return &DatabaseMigrator{DB: db}, nil
}

func (m *DatabaseMigrator) Close() error {
    return m.DB.Close()
}

func (m *DatabaseMigrator) createMigrationsTable() error {
    query := `
    CREATE TABLE IF NOT EXISTS migrations (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        migration TEXT UNIQUE,
        applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`
    _, err := m.DB.Exec(query)
    return err
}

func (m *DatabaseMigrator) getAppliedMigrations() (map[string]bool, error) {
    rows, err := m.DB.Query("SELECT migration FROM migrations")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    applied := make(map[string]bool)
    for rows.Next() {
        var migration string
        if err := rows.Scan(&migration); err != nil {
            return nil, err
        }
        applied[migration] = true
    }
    return applied, nil
}

func (m *DatabaseMigrator) applyMigration(mig migrations.Migration) error {
    tx, err := m.DB.Begin()
    if err != nil {
        return err
    }

    sqlStatements := mig.GetSQL()

    if _, err := tx.Exec(sqlStatements); err != nil {
        tx.Rollback()
        return fmt.Errorf("failed to apply migration %s: %w", mig.Name(), err)
    }

    if _, err := tx.Exec("INSERT INTO migrations(migration) VALUES (?)", mig.Name()); err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit()
}

func (m *DatabaseMigrator) Migrate() error {
    if err := m.createMigrationsTable(); err != nil {
        return err
    }

    applied, err := m.getAppliedMigrations()
    if err != nil {
        return err
    }

    for _, mig := range migrations.Migrations {
        if applied[mig.Name()] {
            logging.Printf("Skipping migration %s (already applied)", mig.Name())
            continue
        }

        logging.Printf("Applying migration %s", mig.Name())
        if err := m.applyMigration(mig); err != nil {
            return err
        }
    }

    return nil
}