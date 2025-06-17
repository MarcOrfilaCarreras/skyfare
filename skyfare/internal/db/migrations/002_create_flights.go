package migrations

type migration002 struct{}

func (m *migration002) Name() string {
    return "002_create_flights"
}

func (m *migration002) GetSQL() string {
    return `
    CREATE TABLE IF NOT EXISTS flights (
        id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
        origin TEXT NOT NULL,
        destination TEXT NOT NULL,
        departure_date DATETIME NOT NULL,
        price REAL NOT NULL,
        promotion BOOLEAN NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (origin) REFERENCES airports(code),
        FOREIGN KEY (destination) REFERENCES airports(code),
        UNIQUE(origin, destination, departure_date)
    );
    
    CREATE TRIGGER IF NOT EXISTS update_flights_updated_at
    AFTER UPDATE ON flights
    FOR EACH ROW
    BEGIN
        UPDATE flights SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
    END;
    `
}

func init() {
    Migrations = append(Migrations, &migration002{})
}
