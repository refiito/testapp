package main

var schema = `
CREATE TABLE IF NOT EXISTS users(
  user_id SERIAL PRIMARY KEY,
  email TEXT NOT NULL,
  password TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX index_user_email_unique ON users USING btree(email);`

func runDBCheck() {
	if config.DB == nil {
		panic("Configuration not yet initialized")
	}
	config.DB.MustExec(schema)
}
