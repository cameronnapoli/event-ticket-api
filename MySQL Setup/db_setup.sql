CREATE DATABASE `Tickets`;

CREATE USER 'admin'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON Tickets.* TO 'admin'@'localhost' IDENTIFIED BY 'password';

CREATE TABLE `Tickets`.`SoldTickets` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `ticketType` varchar(20) NOT NULL,
  `email` varchar(100) NOT NULL,
  `lastUpdateTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
)