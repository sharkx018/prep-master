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
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('done', 'pending', 'in-progress')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_items_status ON items(status);
CREATE INDEX IF NOT EXISTS idx_items_category ON items(category);
CREATE INDEX IF NOT EXISTS idx_items_category_status ON items(category, status);
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
    -- First, temporarily drop the check constraint if it exists
    ALTER TABLE items DROP CONSTRAINT IF EXISTS items_status_check;
    
    -- Update any 'not-done' values to 'pending'
    UPDATE items SET status = 'pending' WHERE status = 'not-done' OR status NOT IN ('done', 'pending', 'in-progress');
    
    -- Drop and recreate the default constraint
    ALTER TABLE items ALTER COLUMN status DROP DEFAULT;
    ALTER TABLE items ALTER COLUMN status SET DEFAULT 'pending';
    
    -- Add the check constraint back with the correct values
    ALTER TABLE items ADD CONSTRAINT items_status_check CHECK (status IN ('done', 'pending', 'in-progress'));
END $$;
`

const addStarredColumn = `
DO $$ 
BEGIN 
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns 
                   WHERE table_name='items' AND column_name='starred') THEN
        ALTER TABLE items ADD COLUMN starred BOOLEAN NOT NULL DEFAULT false;
        CREATE INDEX IF NOT EXISTS idx_items_starred ON items(starred);
    END IF;
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
