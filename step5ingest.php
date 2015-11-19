<?php
/**
 * This script read all domains from the database and from Alexa top 1M list and
 * - updates database with new rank for existing domains
 * - inserts new domains with rank
 * - removes alexa rank for domains no longer on Alexa list
 */
require_once(__DIR__ . '/config.php');

$result_temp_folder = $argv[1];
$icon_archive = $argv[2];

$STDIN = fopen("php://stdin", 'r');

if ($db->query('START TRANSACTION') === FALSE)
{
	throw new Exception("Can't execute query: ".$db->error);
}

$insert_stmt = $db->prepare('INSERT INTO icons (domain_id, retrieved, checksum) VALUES (?, ?, ?)');
if (!$insert_stmt) {
	throw new Exception("Can't prepare statement: ".$db->error);
}

$domain_id = 0;
$retrieved = "";
$checksum = "";

if (!$insert_stmt->bind_param('iss', $domain_id, $retrieved, $checksum)) {
	throw new Exception("Can't bind params: ".$insert_stmt->error);
}

$update_stmt = $db->prepare('UPDATE domains SET last_fetch_time = ?, last_hash = ? WHERE id = ?');
if (!$update_stmt) {
	throw new Exception("Can't prepare statement: ".$db->error);
}

if (!$update_stmt->bind_param('ssi', $retrieved, $checksum, $domain_id)) {
	throw new Exception("Can't bind params: ".$update_stmt->error);
}

$counter = 0;
while($line = fgets($STDIN)) {
	$line = rtrim(trim($line), ',');
	$icon_data = json_decode($line, TRUE);

	/**
	 * If icon was not changed for whatever reason, skip it
	 */
	if (!$icon_data['changed']) {
		continue;
	}

	$domain = $icon_data['domain'];
	$domain_id = $icon_data['id'];
	$retrieved = $icon_data['new_fetch_time'];
	$checksum = $icon_data['new_hash'];

	/**
	 * Let's first move file to a new, permanent location
	 */
	
	// two-level cache dir based on domain name's md5
	$domain_hash = md5($icon_data['domain']);

	if (!file_exists($result_temp_folder . "/" . $domain . ".png")) {
		continue;
	}

	$icon_path = "$icon_archive/" . substr($domain_hash, 0, 2) . "/" . substr($domain_hash, 2, 2) . "/";
	@mkdir($icon_path, 0755, TRUE);
	copy($result_temp_folder . "/" . $domain . ".png", $icon_path . "/" . $domain . "-" . $checksum . ".png");

	/**
	 * Now, let's insert new icon entry into database
	 */
	if (!$insert_stmt->execute()) {
		throw new Exception("Can't execute statement: ".$insert_stmt->error);
	}
	
	/**
	 * And update domain entry with new hash and retrieval date
	 */
	if (!$update_stmt->execute()) {
		throw new Exception("Can't execute statement: ".$update_stmt->error);
	}

	$counter++;
}

$insert_stmt->close();

if ($db->query('COMMIT') === FALSE)
{
	throw new Exception("Can't execute query: ".$db->error);
}

fclose($STDIN);

echo "Inserted $counter new icons\n";
