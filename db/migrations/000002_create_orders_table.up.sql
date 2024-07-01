CREATE TABLE IF NOT EXISTS orders (
    Id              UUID,
    OrderNum        varchar(255) NOT NULL, 
	Login           varchar(255) NOT NULL,    
	CreatedAt       timestamp DEFAULT (now()), 
    Status          varchar(255) DEFAULT 'NEW',		
    AddedBonuses    double precision ,
	SpendedBonuses  double precision DEFAULT 0, 
    SpendDate       timestamp,
PRIMARY KEY (ID));

