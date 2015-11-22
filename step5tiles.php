<?php
/**
 * This script reads all domains with existing icons ordered by their rank and outputs tile creation jobs.
 *
 * Then outputs jobs in JSON format:
 * {"tile":"1","icons":[{"domain":"google.com","hash":"e2020bf4f2b65f62434c62ea967973140b3300df"},{"domain":"youtube.com","hash":"c64bd4efc069a6780274dacbb566cb5d9717c08c"}]}
 */
require_once(__DIR__ . '/config.php');

$tiles_folder = $argv[1];
$html_file = $argv[2];

if ($stmt = $db->prepare("SELECT domain, last_hash FROM domains WHERE last_hash IS NOT NULL ORDER BY alexa_rank"))
{
	if (!$stmt->execute())
	{
		throw new Exception("Can't execute statement: ".$stmt->error);
	}

	if (!$stmt->bind_result($domain, $hash))
	{
		throw new Exception("Can't bind result: ".$stmt->error);
	}

	$tile_size = 20 * 20; // total icons in one tile

	$tile_number = 1;

	$num_icons = 0;
	$icons = array();

	$HTML = fopen($html_file, "w");
	while($stmt->fetch()) {
		$icons[] = array($domain, $hash);

		$num_icons++;

		if ($num_icons >= $tile_size) {
			$TILE = fopen($tiles_folder . "/tile_$tile_number.json", "w");
			fwrite($TILE, json_encode(array(
				'tile' => $tile_number,
				'icons' => $icons
			)));
			fclose($TILE);

			fwrite($HTML, '<img src="tile_' . $tile_number . '.png" data-tile="' . $tile_number . '">' . "\n");

			$tile_number++;

			// resetting the batch
			$num_icons = 0;
			unset($icons);
			$icons = array();			
		}
	}
	fclose($HTML);

	$stmt->close();
}
else
{
	throw new Exception("Can't prepare statement: ".$db->error);
}
