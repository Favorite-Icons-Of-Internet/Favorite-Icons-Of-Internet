#!/bin/perl

use strict;

use Image::Magick;
use Data::Dumper;
use POSIX qw(ceil floor);
use Getopt::Attribute;
use WWW::Favicon;
use JSON;
use Digest::MD5;

our $fetch : Getopt(fetch!);
if (!defined($fetch)) {	$fetch = 1 }

our $genimages : Getopt(genimages!);
if (!defined($genimages)) { $genimages = 1 }

our $ignoreproblems : Getopt(ignoreproblems);
if (!defined($ignoreproblems)) { $ignoreproblems = 0 }

our $genpage: Getopt(genpage!);
if (!defined($genpage)) { $genpage = 1 }

our $desired_page_width : Getopt(width=i);
$desired_page_width ||= 1024;

my $favicon = WWW::Favicon->new;

my @icons;

while(<>) {
	chomp;
	my ($num, $domain) = split(/\t|,/);

	$domain=~s/\s//g;

	print $domain."\t";

	if ($genimages) {
		print "images";

		if (!$ignoreproblems && -f "problems/$domain") {
			print " had problems before, skipping\n";
			next;
		}

		if (-f "icons/$domain.ico" && -z "icons/$domain.ico") {
			unlink "icons/$domain.ico";
			print " empty, deleted.";
		}

		if (!-f "icons/$domain.ico" && $fetch) {
			print " checking the page for links to icons... ";

			# default to /favicon.ico
			my $favicon_url = "http://www.$domain/favicon.ico";

			# let's try getting it from the document
			eval {
				$favicon_url = $favicon->detect("http://www.$domain");
				1;
			};

			print " getting $favicon_url ... ";
			`wget "$favicon_url" -O icons/$domain.ico --dns-timeout=5 --read-timeout=30 --connect-timeout=3 -t 1`;
		}

		if (-f "icons/$domain.ico" && -z "icons/$domain.ico") {
			unlink "icons/$domain.ico";
			print " still empty, deleted.";
		}

		if (-f "icons/$domain.ico" && !-z "icons/$domain.ico" && !-f "pngs/$domain.png") {
			print " got it";
			my $image = Image::Magick->new;
			my $x = $image->Read("icons/$domain.ico");
			warn "$x" if "$x";

			foreach my $i (@$image)
			{
				if ($i->Get('width') == 16 && $i->Get('height') == 16)
				{
					$x = $i->Write("pngs/$domain.png");
					warn "$x" if "$x";

					print " saved png";
					last;
				}
			}
		} elsif (-z "icons/$domain.ico"){
			print " don't have it.";

		} elsif (-f "pngs/$domain.png") {
			print " already have png";
		}
	}

	if (-f "pngs/$domain.png") {
		push(@icons, [$num, $domain]);
	} else {
		`touch problems/$domain`;
	}

	print "\n";
}

if (!$genpage) {
	print "Done downloading icons. Exiting.\n";
	exit;
}

my $chunk = 1;

my $tile_width = 18; # should be more then 16
my $tile_height = 18; # should be more then 16

my $scrollbar_width = 20;

my $margin = 1; # margin between image and border

my $geometry = '16x16+'.($tile_width-16-$margin).'+'.($tile_height-16-$margin);

my $tiles_x = floor(($desired_page_width - $scrollbar_width) / $tile_width);
my $tiles_y_max = 100;

my $page_width = $tiles_x * $tile_width;

my $totalitems = floor($#icons / $tiles_x) * $tiles_x; # to make sure last line is full
@icons = splice(@icons, 0, $totalitems);

my $fullheight = $tile_height * $totalitems / $tiles_x;

#print "\ntile_width: $tile_width\ntile_height: $tile_height\nmargin: $margin\ngeometry: $geometry\ndesired_page_width: $desired_page_width\npage_width: $page_width\ntiles_x: $tiles_x\ntiles_y_max: $tiles_y_max\ntotalitems: $totalitems\nfullheight: $fullheight\n\n";

#my $html = '<div style="height: '.$fullheight.'px">'."\n";
my $html = '';

print "Generating image: ";

my @icons2 = @icons; # copying the list for imagemap generation

#my @hosts = ('http://www.favoriteiconsofinternet.com/', 'http://208.109.208.18/favoriteiconsofinternet.com/', 'http://favoriteiconsofinternet.com/');
#my @hosts = ('http://208.109.208.18/favoriteiconsofinternet.com/', 'http://d2kw54zur50eub.cloudfront.net/');
#my @hosts = ('http://favoriteiconsofinternet.net/', 'http://d2kw54zur50eub.cloudfront.net/');
#my @hosts = ('http://favoriteiconsofinternet.net/', 'http://favoriteiconsofinternet.org/', 'http://d2kw54zur50eub.cloudfront.net/');
#my @hosts = ('http://www.favoriteiconsofinternet.net/', 'http://www.favoriteiconsofinternet.org/', 'http://favoriteiconsofinternet.net/', 'http://favoriteiconsofinternet.org/');
#my @hosts = ('http://208.109.208.18/favoriteiconsofinternet.com/', 'http://favoriteiconsofinternet.com/');
#my @hosts = ('http://favoriteiconsofinternet.net/', 'http://favoriteiconsofinternet.org/');
my @hosts = ('http://www.favoriteiconsofinternet.com/');

my $conn_per_host = 2;
my $current_host = 0;
my $request_counter = 1;

while(my @subset = splice(@icons, 0, $tiles_x * $tiles_y_max))
{
	my $large;

	if ($genimages) {
		$large = Image::Magick->new;
		my @pngs= map { "pngs/".$_->[1].".png"; } @subset;
		my $err = $large->Read(@pngs);
		warn $err if $err;
	}

	my $outfile = 'large_'.$chunk.'_'.$desired_page_width.'x'.$tiles_y_max.'.png';

	if ($genimages) {
		my $result = $large->Montage(
			background => 'white',
			geometry => $geometry,
			tile => $tiles_x.'x'.ceil($#subset / $tiles_x),
			label => undef,
			title => undef
		);
		$result->Write("PNG8:$outfile");
	}

	# md5
	open(FILE, $outfile) or die "Can't open '$outfile': $!";
	binmode(FILE);
	my $hash = Digest::MD5->new->addfile(*FILE)->hexdigest;

	$html .= '<img src="'.$hosts[$current_host].'hash/'.$hash.'/'.$outfile.'" alt="Favorite icons of internet" width="'.$page_width.'" height="'.(int(ceil($#subset / $tiles_x))*$tile_height + 2).'" usemap="#chunk'.$chunk.'" />'."\n";

	$request_counter++;

	if ($request_counter > $conn_per_host-1) {
		$current_host++;
		$request_counter = 0;
	}

	if ($current_host > $#hosts) {
		$current_host = 0;
	}

	$chunk++;

	$large = undef;

	print ".";
}

#$html .= "</div>\n";

print "\nGenerating imagemaps:\n";

$chunk = 1;
#my $sitemap = '';

my $jsoncalls = '';

while(my @subset = splice(@icons2, 0, $tiles_x * $tiles_y_max))
{
	print "Chunk: $chunk\n";

#	$sitemap .= '<map name="chunk'.$chunk.'">'."\n";

	my $x = 0;
	my $y = 0;

	my @chunkdomains = ();

	foreach my $entry (@subset) {
#		$sitemap .= '<area href="http://www.'.$entry->[1].'/" title="'.$entry->[1].'" shape="rect" coords="'.($x + 1).','.($y + 1).','.($x + 18).','.($y + 18).'" />'."\n";

		push(@chunkdomains, $entry->[1]);

		$x += $tile_width;
		print ".";

		if ($x >= $page_width) {
			$x = 0;
			$y += $tile_height;

			print "\n";
		}
	}
#	$sitemap .= '</map>'."\n\n";

	my $body = to_json(\@chunkdomains, {utf8 => 1, pretty => 1, indent_length => 0});

	open(JS, '>chunk'.$chunk.'.js');
	print JS 'createImageMaps('.$chunk.', '.$body.');';
	close(JS);

	# md5
	my $hash = Digest::MD5->new->add($body)->hexdigest;

	$jsoncalls .= '<script src="'.$hosts[$current_host].'hash/'.$hash.'/chunk'.$chunk.'.js" async="async"></script>'."\n";

	$request_counter++;

	if ($request_counter > $conn_per_host-1) {
		$current_host++;
		$request_counter = 0;
	}

	if ($current_host > $#hosts) {
		$current_host = 0;
	}

	$chunk++;
	print "\n";
}
open(TPL, 'index.tpl');
open(HTML, '>index.html');
while(<TPL>){
	s/######ICONS######/$html/;
#	s/######MAP######/$sitemap/;
	s/######JSONCALLS######/$jsoncalls/;
	s/######PAGEWIDTH######/$page_width/;
	print HTML;
}
close(TPL);
close(HTML);
