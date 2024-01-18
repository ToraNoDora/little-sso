package postgres

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
)

const (
	usersTable       = "users"
	permissionsTable = "users_permissions"

	adminsTable      = "admins"
	groupsTable      = "groups"
	groupsRolesTable = "groups_roles"
	rolesTable       = "roles"

	appsTable = "apps"
)

func NewConfig(h, p, u, ps, dbn, sslm string) *Config {
	return &Config{
		Host:     h,
		Port:     p,
		Username: u,
		Password: ps,
		DBName:   dbn,
		SSLMode:  sslm,
	}
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open(
		"postgres",
		fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
		),
	)

	if err != nil {
		log.Fatalf("Error open postgres db: %s", err.Error())
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Error ping postgres db: %s", err.Error())
		return nil, err
	}

	return db, nil
}
