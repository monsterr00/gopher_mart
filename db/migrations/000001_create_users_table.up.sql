CREATE TABLE IF NOT EXISTS users (
    Id             UUID,
	Login          varchar(255) NOT NULL,    
	Password       varchar(255) NOT NULL,    
	CreatedAt      timestamp DEFAULT (now()), 
	Balance        double precision DEFAULT 0,   
	SpendedBonuses double precision DEFAULT 0, 
PRIMARY KEY (ID));