<?php
/**
 * This script read all domains from the database and from Alexa top 1M list and
 * - updates database with new rank for existing domains
 * - inserts new domains with rank
 * - removes alexa rank for domains no longer on Alexa list
 */
require_once(__DIR__ . '/config.php');

define('BATCH_INSERT_SIZE', 2000);

$blacklist = array_map(function($line) {
	return rtrim($line);
}, file('blacklist.txt'));

function fix_domain($domain) {

	return $domain;
}

function finish_db($db) {
	if ($db->query('COMMIT') === FALSE)
	{
		throw new Exception("Can't execute query: ".$db->error);
	}
}

if ($db->query('START TRANSACTION') === FALSE)
{
	throw new Exception("Can't execute query: ".$db->error);
}

// if table has no entries, insert one with default value
if ($db->query('UPDATE domains SET alexa_rank = NULL') === FALSE)
{
	throw new Exception("Can't execute query: ".$db->error);
}

$rank = 0;
$domain = '';

$STDIN = fopen("php://stdin", 'r');

$line = TRUE;
while($line !== FALSE) {
	$batch_insert = "INSERT INTO domains (alexa_rank, domain) VALUES ";
	$first = TRUE;
	for ($i = 0; $i < BATCH_INSERT_SIZE; $i++) {
		$line = fgets($STDIN);
		if ($line == FALSE) {
			break 1;
		}

		$pair = explode(',', $line);

		// something is wrong with this line, exiting.
		if (count($pair) !== 2) {
			echo "[$domain]";
			finish_db($db);
			exit;
		}

		$rank = $pair[0];
		$domain = trim($pair[1]);

		// if domain is invalid, skip it
		if (strpos($domain, '/') !== FALSE) {
			continue;
		}

		if (in_array($domain, $blacklist)) {
			continue;
		}

		if ($first) {
			$first = FALSE;
		} else {
			$batch_insert .= ",";
		}

		$batch_insert .= "(" . intval($rank) . ",'" . $db->real_escape_string($domain) . "')"; 

	}
	$batch_insert .= " ON DUPLICATE KEY UPDATE alexa_rank = VALUES(alexa_rank)";

	if ($i == 0) {
		continue;
	}

	if (!$db->query($batch_insert))
	{
		throw new Exception("Can't execute statement: ".$stmt->error);
	}
}

fclose($STDIN);

finish_db($db);
