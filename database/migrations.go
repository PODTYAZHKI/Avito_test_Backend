package database

import (
	"github.com/jmoiron/sqlx"
	"log"
)

// Migrate выполняет миграции для создания необходимых таблиц
func Migrate(db *sqlx.DB) {
	schema := `
	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
	        CREATE TYPE tender_status AS ENUM ('Created', 'Published', 'Closed');
	END IF;
	END$$;
	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_service_type') THEN
	        CREATE TYPE tender_service_type AS ENUM ('Construction', 'Delivery', 'Manufacture');
	END IF;
	END$$;

	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_author_type') THEN
	        CREATE TYPE bid_author_type AS ENUM ('Organization', 'User');
	END IF;
	END$$;

	DO $$
	BEGIN
	    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'bid_status') THEN
	        CREATE TYPE bid_status AS ENUM ('Created', 'Published', 'Canceled', 'Approved', 'Rejected');
	END IF;
	END$$;


    CREATE TABLE IF NOT EXISTS tender (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        name VARCHAR(100) NOT NULL,
        description TEXT,
        service_type tender_service_type,
        status tender_status,
        version INT,
        organization_id uuid REFERENCES organization(id) ON DELETE CASCADE,
		user_id uuid REFERENCES employee(id) ON DELETE CASCADE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );

	CREATE TABLE IF NOT EXISTS tender_versions (
    	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    	tender_id uuid REFERENCES tender(id) ON DELETE CASCADE,
    	name VARCHAR(100),
    	description TEXT,
    	service_type tender_service_type,
    	status tender_status,
    	version INT,
		organization_id uuid REFERENCES organization(id) ON DELETE CASCADE,
		user_id uuid REFERENCES employee(id) ON DELETE CASCADE,
    	created_at TIMESTAMP,
    	updated_at TIMESTAMP
		);

	CREATE TABLE IF NOT EXISTS bid_versions (
    	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    	bid_id uuid REFERENCES bid(id) ON DELETE CASCADE,
    	name VARCHAR(100),
    	description TEXT,
		status bid_status,
		tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
		author_id UUID,
		author_type bid_author_type,
		version INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

	CREATE TABLE IF NOT EXISTS bid (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		name VARCHAR(100),
    	description TEXT,
		status bid_status,
		tender_id UUID REFERENCES tender(id) ON DELETE CASCADE,
		author_id UUID,
		author_type bid_author_type,
		version INT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
     `
	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}
}

// DROP TABLE IF EXISTS tender;
// DO $$
// BEGIN
//     IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_status') THEN
//         EXECUTE 'DROP TYPE tender_status';
//     END IF;

//     IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'tender_service_type') THEN
//         EXECUTE 'DROP TYPE tender_service_type';
//     END IF;
// END$$;
