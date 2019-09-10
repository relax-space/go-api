-- --------------------------------------------------------
-- 主机:                           10.173.208.142
-- 服务器版本:                        5.7.25-log - MySQL Community Server (GPL)
-- 服务器操作系统:                      Linux
-- HeidiSQL 版本:                  9.5.0.5295
-- --------------------------------------------------------

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET NAMES utf8 */;
/*!50503 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;


-- 导出 fruit 的数据库结构
CREATE DATABASE IF NOT EXISTS `fruit` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_bin */;
USE `fruit`;

-- 导出  表 fruit.store 结构
DROP TABLE IF EXISTS `store`;
CREATE TABLE IF NOT EXISTS `store` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `code` varchar(255) DEFAULT NULL,
  `name` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7119 DEFAULT CHARSET=utf8;

-- 正在导出表  fruit.store 的数据：~7,099 rows (大约)
/*!40000 ALTER TABLE `store` DISABLE KEYS */;
INSERT INTO `store` (`id`, `code`, `name`) VALUES
	(1, '10001', '果想你'),
	(2, '10002', '果为鲜');
/*!40000 ALTER TABLE `store` ENABLE KEYS */;

/*!40101 SET SQL_MODE=IFNULL(@OLD_SQL_MODE, '') */;
/*!40014 SET FOREIGN_KEY_CHECKS=IF(@OLD_FOREIGN_KEY_CHECKS IS NULL, 1, @OLD_FOREIGN_KEY_CHECKS) */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
