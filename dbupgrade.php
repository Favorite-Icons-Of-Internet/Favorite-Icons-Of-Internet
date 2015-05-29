<?php
/*
 * Copy this script to the folder above and populate $versions array with your migrations
 * For more info see: http://www.dbupgrade.org/Main_Page#Migrations_($versions_array)
 *
 * Note: this script should be versioned in your code repository so it always reflects current code's
 *       requirements for the database structure.
*/
require_once(__DIR__ . '/DBUpgrade/lib.php');

$versions = array();

// Add new migrations on top, right below this line.
/* -------------------------------------------------------------------------------------------------------
 * VERSION _
 * ... add version description here ...
*/
/*
$versions[_]['up'][] = "";
$versions[_]['down'][]	= "";
*/

/* -------------------------------------------------------------------------------------------------------
 * VERSION _
 * Getting rid of the icon blob (will store on file system
*/
$versions[2]['up'][] = "ALTER TABLE `icons` DROP `optimized_png_content`";
$versions[2]['down'][]	= "ALTER TABLE `icons` ADD `optimized_png_content` BLOB NOT NULL COMMENT 'Content of optimized image file converted to PNG'";

/* -------------------------------------------------------------------------------------------------------
 * VERSION 1
 * First draft of schema
*/
$versions[1]['up'][] = "CREATE TABLE `domains` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Domain ID',
`alexa_rank` int(10) unsigned NULL DEFAULT NULL COMMENT 'Rank of the domain, e.g. global order',
`domain` varchar(100) NOT NULL COMMENT 'Domain name',
PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Domains to collect icons for' AUTO_INCREMENT=1
";

$versions[1]['up'][] = "CREATE TABLE `icons` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Individual icon ID',
`domain_id` int(10) unsigned NOT NULL COMMENT 'Foreign key column linking icon to a domain',
`retrieved` datetime DEFAULT NULL COMMENT 'Date of icon retrieval',
`checksum` bigint(20) unsigned NOT NULL COMMENT 'Checksum of original icon file',
`average_color` bigint(6) unsigned DEFAULT NULL COMMENT 'Integer representation of average color for the icon',
`optimized_png_content` blob NOT NULL COMMENT 'Content of optimized image file converted to PNG',
PRIMARY KEY (`id`),
CONSTRAINT `icon_domain` FOREIGN KEY (`domain_id`) REFERENCES `domains` (`id`)
ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 AUTO_INCREMENT=1
";

$versions[1]['down'][] = "DROP TABLE icons";
$versions[1]['down'][] = "DROP TABLE domains";

require_once(__DIR__ . '/config.php');
// creating DBUpgrade object with your database credentials and $versions defined above
$dbupgrade = new DBUpgrade($db,	$versions);
require_once(__DIR__ . '/DBUpgrade/client.php');
