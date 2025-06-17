package migrations

type migration001 struct{}

func (m *migration001) Name() string {
    return "001_create_airports"
}

func (m *migration001) GetSQL() string {
    return `
    CREATE TABLE airports (
        id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
        code TEXT UNIQUE NOT NULL,
        name TEXT,
        country TEXT,
        lat REAL,
        lng REAL
    );
    `
}

func init() {
    Migrations = append(Migrations, &migration001{})
}
