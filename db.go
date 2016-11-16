package main

import (
	"errors"
	"github.com/jameskeane/bcrypt"
)

var schema = `
CREATE TABLE IF NOT EXISTS users(
  user_id SERIAL PRIMARY KEY,
  email TEXT NOT NULL,
  password TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
BEGIN
  IF to_regclass('public.index_user_email_unique') IS NULL THEN
    CREATE UNIQUE INDEX index_user_email_unique ON users USING btree(email);
  END IF;
END$$;`

var insertUser = `INSERT INTO users (email, password) VALUES ($1, $2)`
var selectUser = `SELECT user_id, email, password FROM users WHERE lower(email) = lower($1)`

func runDBCheck() {
	if config.DB == nil {
		panic("Configuration not yet initialized")
	}
	config.DB.MustExec(schema)
}

func encryptPassword(pwd string) (string, error) {
	if len(pwd) < 6 {
		return "", errors.New("Password should be at least 6 characters")
	}

	salt, err := bcrypt.Salt(10)
	if err != nil {
		return "", err
	}

	return bcrypt.Hash(pwd, salt)
}

func createUser(email, password string) (err error) {
	existingUser, err := getUser(email)
	if existingUser.ID > 0 {
		err = errors.New("E-mail address already in use, user exists")
	}
	if err != nil {
		return
	}

	cryptedPassword, err := encryptPassword(password)
	if err != nil {
		return
	}

	_, err = config.DB.Exec(insertUser, email, cryptedPassword)
	return
}

func getUser(email string) (usr User, err error) {
	err = config.DB.Get(&usr, selectUser, email)
	return
}

func (usr *User) authenticate(password string) (err error) {
	if !bcrypt.Match(password, usr.CryptedPassword) {
		err = errors.New("Password doesn't match")
	}
	return
}
