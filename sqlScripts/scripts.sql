CREATE TABLE Users (
	Id SERIAL PRIMARY KEY,
	Email varchar(50) UNIQUE NOT NULL,
	PasswordHash varchar(250) NOT NULL,
	IsVerified int NOT NULL,
	IsAdmin int NOT NULL
)

----------

CREATE TABLE Demoscenes (
	Id SERIAL PRIMARY KEY,
	Name varchar(50) UNIQUE NOT NULL,
	RepositoryName varchar(250) NOT NULL,
	BranchName varchar(250) NOT NULL,
	IsDeleted int NOT NULL
)