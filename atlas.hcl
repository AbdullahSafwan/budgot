env "local" {
  url = "sqlite://budgot.db?_fk=1"

  migration {
    dir = "file://internal/ent/migrate/migrations"
  }
}

env "prod" {
  url = "sqlite://budgot.db?_fk=1"
  migration {
    dir = "file://internal/ent/migrate/migrations"
  }
}
