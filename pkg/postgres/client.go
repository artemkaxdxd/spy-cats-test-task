package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

type (
	Database interface {
		// Instance returns the GORM DB instance.
		Instance() *gorm.DB
		// Close closes the underlying database connection.
		Close() error
		// Ping checks connectivity to the database.
		Ping(ctx context.Context) error
	}

	Postgres struct {
		db    *gorm.DB
		sqlDB *sql.DB
	}
)

// New creates a new Postgres instance using the provided parameters.
func New(
	ctx context.Context,
	user, password, host, port, dbname string,
	isDev bool,
) (*Postgres, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host, user, password, dbname, port,
	)

	logLvl := gormLogger.Error
	if isDev {
		logLvl = gormLogger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(logLvl),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(0)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctxTimeout); err != nil {
		return nil, err
	}

	return &Postgres{db: db, sqlDB: sqlDB}, nil
}

// Instance returns the GORM DB instance.
func (p *Postgres) Instance() *gorm.DB {
	return p.db
}

// Close closes the underlying database connection.
func (p *Postgres) Close() error {
	if p.sqlDB == nil {
		return nil
	}
	return p.sqlDB.Close()
}

// Ping checks connectivity to the database.
func (p *Postgres) Ping(ctx context.Context) error {
	if p.sqlDB == nil {
		return fmt.Errorf("database connection is closed")
	}
	return p.sqlDB.PingContext(ctx)
}
