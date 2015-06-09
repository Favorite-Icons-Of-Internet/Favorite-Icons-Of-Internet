<?php
/**
 * This script read all domains and corresponding icon metadata from the database that need to be fetched,
 * e.g. those that have Alexa rankings
 *
 * Then outputs initial jobs in semi-JSON format:
 * {"id": "1","domain": "google.com","previous_hash": "aqwerdqweafv45vacrqe","previous_fetch_time": "2015-04-23T18:25:43.511Z"},
 */
require_once(__DIR__ . '/config.php');

if ($db->query("SET time_zone = '+00:00'") === FALSE) {
	throw new Exception("Can't set UTC timezone: ".$db->error);
}

if ($stmt = $db->prepare("SELECT id, domain, last_hash, DATE_FORMAT(last_fetch_time,'%Y-%m-%dT%TZ') FROM domains WHERE alexa_rank IS NOT NULL"))
{
	if (!$stmt->execute())
	{
		throw new Exception("Can't execute statement: ".$stmt->error);
	}

	if (!$stmt->bind_result($id, $domain, $previous_hash, $previous_fetch_time))
	{
		throw new Exception("Can't bind result: ".$stmt->error);
	}

	while($stmt->fetch()) {
		if (is_null($previous_hash)) {
			echo json_encode(array(
				'id' => $id,
				'domain' => $domain
			));
		} else {
			echo json_encode(array(
				'id' => $id,
				'domain' => $domain,
				'previous_hash' => $previous_hash,
				'previous_fetch_time' => $previous_fetch_time
			));
		}
		echo ",\n";
	}

	$stmt->close();
}
else
{
	throw new Exception("Can't prepare statement: ".$db->error);
}
