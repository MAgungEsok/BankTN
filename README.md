# BankTN

Sql with the database name "bank"

CREATE TABLE `costumer` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `acountNumber` varchar(13) NOT NULL,
  `email` varchar(20) NOT NULL,
   PRIMARY KEY (id)
);

CREATE TABLE `transaction` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `idCostumer` int(10) NOT NULL,
  `deposit` bigint(20) NOT NULL,
  PRIMARY KEY (id)
);


ps:
I made it with MySQL

Design a software architecture, can be check in the image with the name "use case diagram" and "class diagram".

Thanks.
