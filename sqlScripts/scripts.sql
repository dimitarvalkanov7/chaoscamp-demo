CREATE TABLE Users (
	Id SERIAL PRIMARY KEY,
	Email varchar(50) UNIQUE NOT NULL,
	PasswordHash varchar(250) NOT NULL
)