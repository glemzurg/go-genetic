# ************************************************************
# Sequel Pro SQL dump
# Version 4096
#
# http://www.sequelpro.com/
# http://code.google.com/p/sequel-pro/
#
# Host: 127.0.0.1 (MySQL 5.6.21)
# Database: genetic
# Generation Time: 2015-05-16 22:18:15 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table experiment
# ------------------------------------------------------------

DROP TABLE IF EXISTS `experiment`;

CREATE TABLE `experiment` (
  `experimentid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `experiment` varchar(32) NOT NULL DEFAULT '',
  `datetime` datetime NOT NULL,
  `config` blob NOT NULL,
  `scorer` blob NOT NULL,
  `selector_type` varchar(128) NOT NULL,
  `selector` blob NOT NULL,
  PRIMARY KEY (`experimentid`),
  UNIQUE KEY `experiment` (`experiment`,`datetime`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table experiment_end
# ------------------------------------------------------------

DROP TABLE IF EXISTS `experiment_end`;

CREATE TABLE `experiment_end` (
  `experimentid` int(11) unsigned NOT NULL,
  `end_reason` varchar(256) NOT NULL DEFAULT '',
  `datetime` datetime NOT NULL,
  `generation_num` bigint(20) unsigned NOT NULL,
  `results` longblob NOT NULL,
  PRIMARY KEY (`experimentid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table experiment_generation
# ------------------------------------------------------------

DROP TABLE IF EXISTS `experiment_generation`;

CREATE TABLE `experiment_generation` (
  `experimentid` int(11) unsigned NOT NULL,
  `generation_num` bigint(11) unsigned NOT NULL,
  `datetime` datetime NOT NULL,
  `highest_experiment_score` double NOT NULL,
  `stagnant_generations` int(11) NOT NULL,
  `details` blob NOT NULL,
  PRIMARY KEY (`experimentid`,`generation_num`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;



# Dump of table experiment_generation_species
# ------------------------------------------------------------

DROP TABLE IF EXISTS `experiment_generation_species`;

CREATE TABLE `experiment_generation_species` (
  `experimentid` int(11) unsigned NOT NULL,
  `generation_num` bigint(11) unsigned NOT NULL,
  `species_fingerprint` char(32) NOT NULL DEFAULT '',
  `specimens` bigint(20) unsigned NOT NULL,
  `highest_score` float(20,2) NOT NULL,
  `highest_bonus` float(20,2) NOT NULL,
  `highest_species_score` float(20,2) NOT NULL,
  PRIMARY KEY (`experimentid`,`generation_num`,`species_fingerprint`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
