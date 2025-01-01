package migrations

import (
	"github.com/mibrahim2344/identity-service/internal/domain/models"
	"gorm.io/gorm"
)

// Migrate performs database migrations using GORM
func Migrate(db *gorm.DB) error {
	// Create enum types first
	err := db.Exec(`DO $$ BEGIN
		CREATE TYPE user_status AS ENUM ('active', 'inactive', 'pending');
		EXCEPTION WHEN duplicate_object THEN NULL;
	END $$;`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`DO $$ BEGIN
		CREATE TYPE user_role AS ENUM ('admin', 'user');
		EXCEPTION WHEN duplicate_object THEN NULL;
	END $$;`).Error
	if err != nil {
		return err
	}

	// Auto-migrate the User model
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	// Create indexes
	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)").Error
	if err != nil {
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)").Error
	if err != nil {
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_status ON users(status)").Error
	if err != nil {
		return err
	}

	err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at)").Error
	if err != nil {
		return err
	}

	return nil
}

// Rollback removes all migrations
func Rollback(db *gorm.DB) error {
	// Drop the users table
	err := db.Migrator().DropTable(&models.User{})
	if err != nil {
		return err
	}

	// Drop enum types
	err = db.Exec(`DROP TYPE IF EXISTS user_status`).Error
	if err != nil {
		return err
	}

	err = db.Exec(`DROP TYPE IF EXISTS user_role`).Error
	if err != nil {
		return err
	}

	return nil
}
