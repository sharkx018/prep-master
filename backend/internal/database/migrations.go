package database

import (
	"database/sql"
	"fmt"
	"log"
)

// RunMigrations executes all database migrations
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		createItemsTable,
		createAppStatsTable,
		insertInitialAppStats,
		addSubcategoryColumn,
		fixStatusValues,
		addStarredColumn,
		addAttachmentsColumn,
		createUsersTable,
		createUserProgressTable,
		createUserStatsTable,
		createRefreshTokensTable,
		fixUsersUniqueConstraint,
		addUserRoleColumn,
		addUserProgressStarredColumn,
		addUserStatsCompletedAllCountColumn,
	}

	for i, migration := range migrations {
		if err := executeMigration(db, migration); err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// executeMigration executes a single migration
func executeMigration(db *sql.DB, migration string) error {
	_, err := db.Exec(migration)
	return err
}

// Migration SQL statements
const createItemsTable = `
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    link TEXT NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('dsa', 'lld', 'hld')),
    subcategory VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
`

const createAppStatsTable = `
CREATE TABLE IF NOT EXISTS app_stats (
    id INTEGER PRIMARY KEY DEFAULT 1,
    completed_all_count INTEGER DEFAULT 0,
    CONSTRAINT single_row CHECK (id = 1)
);
`

const insertInitialAppStats = `
INSERT INTO app_stats (id, completed_all_count) 
VALUES (1, 0) 
ON CONFLICT (id) DO NOTHING;
`

const addSubcategoryColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='items' AND column_name='subcategory') THEN
        ALTER TABLE items ADD COLUMN subcategory VARCHAR(100) NOT NULL DEFAULT 'other';
        CREATE INDEX IF NOT EXISTS idx_items_subcategory ON items(subcategory);
        CREATE INDEX IF NOT EXISTS idx_items_category_subcategory ON items(category, subcategory);
    END IF;
END $$;
`

const fixStatusValues = `
DO $$
BEGIN
    -- This migration is no longer needed as status column is handled in user_progress table
    -- No operation needed
END $$;
`

const addStarredColumn = `
DO $$ 
BEGIN 
    -- This migration is no longer needed as starred column is handled in user_progress table
    -- No operation needed
END $$;
`

const addAttachmentsColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='items' AND column_name='attachments') THEN
        ALTER TABLE items ADD COLUMN attachments JSONB DEFAULT '{}';
    END IF;
END $$;
`

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    auth_provider VARCHAR(50) NOT NULL CHECK (auth_provider IN ('email', 'google', 'facebook', 'apple')),
    provider_id VARCHAR(255),
    avatar TEXT,
    email_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_provider ON users(auth_provider, provider_id);
CREATE INDEX IF NOT EXISTS idx_users_active ON users(is_active);

-- Create partial unique index for OAuth providers only (when provider_id is not null)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth_unique ON users(auth_provider, provider_id) WHERE provider_id IS NOT NULL;
`

const createUserProgressTable = `
CREATE TABLE IF NOT EXISTS user_progress (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('done', 'pending', 'in-progress')),
    notes TEXT,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, item_id)
);

CREATE INDEX IF NOT EXISTS idx_user_progress_user_id ON user_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_item_id ON user_progress(item_id);
CREATE INDEX IF NOT EXISTS idx_user_progress_status ON user_progress(status);
CREATE INDEX IF NOT EXISTS idx_user_progress_user_status ON user_progress(user_id, status);
`

const createUserStatsTable = `
CREATE TABLE IF NOT EXISTS user_stats (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_items INTEGER DEFAULT 0,
    completed_items INTEGER DEFAULT 0,
    in_progress_items INTEGER DEFAULT 0,
    pending_items INTEGER DEFAULT 0,
    dsa_completed INTEGER DEFAULT 0,
    lld_completed INTEGER DEFAULT 0,
    hld_completed INTEGER DEFAULT 0,
    current_streak INTEGER DEFAULT 0,
    longest_streak INTEGER DEFAULT 0,
    last_activity_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
`

const createRefreshTokensTable = `
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_revoked BOOLEAN DEFAULT false
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);
`

const fixUsersUniqueConstraint = `
-- Drop the existing unique constraint if it exists
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_auth_provider_provider_id_key;

-- Create partial unique index for OAuth providers only (when provider_id is not null)
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_oauth_unique ON users(auth_provider, provider_id) WHERE provider_id IS NOT NULL;
`

const addUserRoleColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='users' AND column_name='role') THEN
        ALTER TABLE users ADD COLUMN role VARCHAR(20) NOT NULL DEFAULT 'user' CHECK (role IN ('user', 'admin'));
        CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
    END IF;
END $$;
`

const addUserProgressStarredColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='user_progress' AND column_name='starred') THEN
        ALTER TABLE user_progress ADD COLUMN starred BOOLEAN NOT NULL DEFAULT false;
        CREATE INDEX IF NOT EXISTS idx_user_progress_starred ON user_progress(starred);
        CREATE INDEX IF NOT EXISTS idx_user_progress_user_starred ON user_progress(user_id, starred);
    END IF;
END $$;
`

const addUserStatsCompletedAllCountColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='user_stats' AND column_name='completed_all_count') THEN
        ALTER TABLE user_stats ADD COLUMN completed_all_count INTEGER NOT NULL DEFAULT 0;
    END IF;
END $$;
`
