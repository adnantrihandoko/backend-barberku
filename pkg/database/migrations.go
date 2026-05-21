package database

import (
	"context"
	"log/slog"
)

func (db *Database) Migrate(ctx context.Context) error {
	slog.Info("running database migrations")

	migrations := []string{
		createUsersTable,
		createBarbersTable,
		createServicesTable,
		createQueuesTable,
		createStoreSettingsTable,
		createFCMTokensTable,
	}

	for _, migration := range migrations {
		if _, err := db.Pool.Exec(ctx, migration); err != nil {
			return err
		}
	}

	slog.Info("migrations completed")
	return nil
}

const createUsersTable = `
CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL,
	phone VARCHAR(50),
	role VARCHAR(50) NOT NULL DEFAULT 'customer',
	pin_hash VARCHAR(255),
	is_active BOOLEAN DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createBarbersTable = `
CREATE TABLE IF NOT EXISTS barbers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	specialty VARCHAR(255),
	is_active BOOLEAN DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createServicesTable = `
CREATE TABLE IF NOT EXISTS services (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name VARCHAR(255) NOT NULL,
	description TEXT,
	price DECIMAL(10, 2) NOT NULL DEFAULT 0,
	duration INTEGER NOT NULL DEFAULT 30,
	is_active BOOLEAN DEFAULT true,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createQueuesTable = `
CREATE TABLE IF NOT EXISTS queues (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	queue_number INTEGER NOT NULL,
	customer_id VARCHAR(255) NOT NULL,
	customer_name VARCHAR(255) NOT NULL,
	barber_id UUID REFERENCES barbers(id),
	service_id UUID NOT NULL REFERENCES services(id),
	service_name VARCHAR(255) NOT NULL,
	status VARCHAR(50) NOT NULL DEFAULT 'waiting',
	position INTEGER DEFAULT 0,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	called_at TIMESTAMP WITH TIME ZONE,
	completed_at TIMESTAMP WITH TIME ZONE,
	rating INTEGER,
	rating_comment TEXT
);

CREATE INDEX IF NOT EXISTS idx_queues_customer_id ON queues(customer_id);
CREATE INDEX IF NOT EXISTS idx_queues_status ON queues(status);
CREATE INDEX IF NOT EXISTS idx_queues_queue_number ON queues(queue_number);
`

const createStoreSettingsTable = `
CREATE TABLE IF NOT EXISTS store_settings (
	id SERIAL PRIMARY KEY,
	open_hour INTEGER NOT NULL DEFAULT 9,
	close_hour INTEGER NOT NULL DEFAULT 21,
	max_queue_size INTEGER NOT NULL DEFAULT 50,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
`

const createFCMTokensTable = `
CREATE TABLE IF NOT EXISTS fcm_tokens (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	customer_id VARCHAR(255) NOT NULL,
	token TEXT NOT NULL,
	platform VARCHAR(50) NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
	UNIQUE(customer_id, token)
);

CREATE INDEX IF NOT EXISTS idx_fcm_tokens_customer_id ON fcm_tokens(customer_id);
`
