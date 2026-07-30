env "local" {
  url = "mysql://${getenv("DB_USER")}:${getenv("DB_PASSWORD")}@${getenv("DB_HOST")}:${getenv("DB_PORT")}/${getenv("DB_NAME")}"

  migration {
    dir = "file://internal/ent/migrate/migrations"
  }
}

env "prod" {
  url = getenv("DATABASE_URL")
  migration {
    dir = "file://internal/ent/migrate/migrations"
  }
}
