/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

DROP TABLE IF EXISTS `person`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `person` (
  `id` int(9) unsigned NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) NOT NULL,
  `last_name` varchar(100) NOT NULL,
  `sex` enum('male','female','other') NOT NULL,
  KEY `id` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;


LOCK TABLES `person` WRITE;
/*!40000 ALTER TABLE `person` DISABLE KEYS */;
INSERT INTO `person` VALUES (1,'Devon','Ortiz','female'),(2,'Nya','Ullrich','other'),(3,'Chaim','Yundt','female'),(4,'Colton','Rutherford','female'),(5,'Hellen','Borer','female'),(6,'Alaina','Corkery','male'),(7,'Mose','Moore','female'),(8,'Leslie','Hessel','female'),(9,'Cyrus','Medhurst','other'),(10,'Bridie','Swaniawski','female'),(11,'Muhammad','Emmerich','other'),(12,'Howell','Stamm','other'),(13,'Maximo','Glover','male'),(14,'Camylle','Jast','female'),(15,'Stephanie','Goldner','female'),(16,'Murphy','Gutmann','male'),(17,'Breanna','Gutkowski','female'),(18,'Megane','Sauer','male'),(19,'Damon','Stokes','male'),(20,'Blaze','Weissnat','other'),(21,'Delfina','Goyette','other'),(22,'Reynold','Okuneva','female'),(23,'Alice','Cole','male'),(24,'Dedric','Fisher','male'),(25,'Haskell','Stokes','male'),(26,'Jayson','Kshlerin','male'),(27,'Jefferey','Adams','male'),(28,'Yoshiko','Williamson','male'),(29,'Jeremie','Stracke','female'),(30,'Benton','Anderson','male'),(31,'Lauren','Wolff','female'),(32,'Israel','Kirlin','other'),(33,'Ines','Ullrich','male'),(34,'Treva','Kertzmann','other'),(35,'Deborah','Braun','other'),(36,'Athena','Bradtke','other'),(37,'Mohamed','Krajcik','other'),(38,'Norris','Durgan','male'),(39,'Lucas','Tromp','female'),(40,'Romaine','Mante','male'),(41,'Tatyana','Bayer','other'),(42,'Osvaldo','Boyer','male'),(43,'Junius','Wiza','male'),(44,'Justen','Emmerich','female'),(45,'Meda','Johns','female'),(46,'Marco','Ratke','male'),(47,'Guadalupe','Mitchell','female'),(48,'Carleton','Schowalter','female'),(49,'Brain','Jast','other'),(50,'Eloy','Lockman','male'),(51,'Edwin','Schulist','other'),(52,'Florence','Gutmann','male'),(53,'Lee','Langworth','female'),(54,'Jodie','Gleichner','male'),(55,'Shannon','Crooks','other'),(56,'Everardo','Quitzon','other'),(57,'Maryse','Hermann','female'),(58,'Travis','Rolfson','female'),(59,'Sierra','Rosenbaum','male'),(60,'Aida','Steuber','female'),(61,'Abby','Bernier','male'),(62,'Karina','Howe','female'),(63,'Dorthy','Ferry','male'),(64,'Hiram','Turcotte','female'),(65,'Alex','Auer','male'),(66,'Marjory','Murray','male'),(67,'Jaylen','Reichel','other'),(68,'Chyna','Lueilwitz','female'),(69,'Gene','Morar','female'),(70,'Felicity','Yundt','female'),(71,'Lane','Schoen','other'),(72,'Torrey','Gleichner','other'),(73,'Jett','Hoppe','male'),(74,'Murphy','Ortiz','other'),(75,'Lafayette','Raynor','male'),(76,'Daniela','Pfeffer','female'),(77,'Ashlee','Schuppe','female'),(78,'Carrie','Bins','male'),(79,'Garth','Schmeler','other'),(80,'Cortez','Kautzer','other'),(81,'Antonietta','Murazik','other'),(82,'Hipolito','Monahan','other'),(83,'Drew','Breitenberg','other'),(84,'Carmela','Ferry','other'),(85,'Brain','Kozey','female'),(86,'Reva','Mayert','female'),(87,'Brock','Haag','male'),(88,'Jacky','Terry','female'),(89,'Noemy','Schinner','female'),(90,'Branson','Mills','female'),(91,'Danial','Beahan','female'),(92,'Myrtis','Hamill','female'),(93,'Elroy','Hessel','female'),(94,'Ollie','Tromp','male'),(95,'Arnaldo','Monahan','female'),(96,'Audrey','Rodriguez','female'),(97,'Leila','Donnelly','other'),(98,'Madalyn','Pollich','male'),(99,'Sherman','Batz','other'),(100,'Maryjane','Vandervort','other');
/*!40000 ALTER TABLE `person` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2023-01-19 12:00:49
