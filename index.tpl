<HTML>
<HEAD>
<TITLE>Favorite Icons of Internet</TITLE>
<style>
body { margin: 0; overflow-x: hidden; background-color: white}
img { border: 0; display: block; margin: 0; padding: 0; }
</style>
<script>
function createImageMaps(chunk, domains) {
	var x = 0;
	var y = 0;

	var tile_width = 18; // should be more then 16
	var tile_height = 18; // should be more then 16
	var page_width = 990;

	var map = '<map name="chunk' + chunk + '">\n';
	for (var i in domains) {
		map += '<area target="_blank" href="http://www.' + domains[i] + '/" title="' + domains[i] + '" shape="rect" coords="' + (x + 1) + ',' + (y + 1) + ',' + (x + 18) + ',' + (y + 18) + '" />\n';
		x += tile_width;

		if (x >= page_width) {
			x = 0;
			y += tile_height;
		}
	}
	map += '</map>';

	var mapel = document.createElement('div');
	mapel.innerHTML = map;

	document.body.appendChild(mapel);
}

</script>
</HEAD>
<BODY>
######ICONS######

<script type="text/javascript">

  var _gaq = _gaq || [];
  _gaq.push(['_setAccount', 'UA-817839-31']);
  _gaq.push(['_trackPageview']);

  (function() {
    var ga = document.createElement('script'); ga.type = 'text/javascript'; ga.async = true;
    ga.src = ('https:' == document.location.protocol ? 'https://ssl' : 'http://www') + '.google-analytics.com/ga.js';
    var s = document.getElementsByTagName('script')[0]; s.parentNode.insertBefore(ga, s);
  })();

</script>

######JSONCALLS######

</BODY></HTML>
