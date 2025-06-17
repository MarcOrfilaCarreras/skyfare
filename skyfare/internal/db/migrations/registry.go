package migrations

type Migration interface {
    Name() string
    GetSQL() string
}

var Migrations []Migration
