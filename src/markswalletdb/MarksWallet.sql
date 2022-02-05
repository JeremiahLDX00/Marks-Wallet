CREATE database markswalletdb;
USE markswalletdb;

#DROP TABLE Transactions;
#DROP TABLE TokenType;

CREATE TABLE TokenType(
TokenTypeID int(5) auto_increment PRIMARY KEY,
TokenTypeName varchar(5));

CREATE Table Transactions(
TransactionID int(3) auto_increment PRIMARY KEY,
StudentID varchar(9) NOT NULL,
ToStudentID varchar (9) NOT NULL,
TokenTypeID int(5),
TransactionType varchar(30),
Amount int(3),
FOREIGN KEY(TokenTypeID) REFERENCES TokenType(TokenTypeID));

INSERT INTO TokenType(TokenTypeName) VALUES("ETI"), ("CM"), ("CSF"), ("DP"), ("PRG1"), ("DB"), 
("ID"), ("OSNF"), ("PRG2"), ("OOAD"), ("WEB"), ("PFD"), ("SDD");

INSERT INTO Transactions(StudentID, ToStudentID, TokenTypeID, TransactionType, Amount)
VALUES('S10198398','S10183726', 2 , 'Marks Entry', 20), 
('S10198398','S10183726', 1 , 'Marks Entry', 30), 
('S10198398','S10183726', 1 , 'Bidding', -20),
('S10198397','S10183726', 1 , 'Marks Entry', 30);

#To list all tokens and transactions without reference to studentID
SELECT * FROM TokenType;
SELECT * FROM Transactions;

#To list all tokens available for student with balance
SELECT StudentID, TokenType.TokenTypeName, SUM(Amount) FROM Transactions 
INNER JOIN TokenType 
ON Transactions.TokenTypeID = TokenType.TokenTypeID 
WHERE StudentID = "S10198398"
GROUP BY StudentID, Transactions.TokenTypeID;

#To search for individual tokens
SELECT TokenTypeID, TokenTypeName
FROM TokenType
WHERE TokenTypeName = 'SDD';

#To list all the transactions
SELECT TransactionID, StudentID, ToStudentID, TokenType.TokenTypeName, TransactionType, Amount 
FROM Transactions
INNER JOIN TokenType
ON Transactions.TokenTypeID = TokenType.TokenTypeID
WHERE StudentID = "S10198398";