#!/bin/perl
#
#
use strict;

use WWW::AppleTouchIcon;
my $apple_touch_icon = WWW::AppleTouchIcon->new;
my $apple_touch_icon_url = $apple_touch_icon->detect('http://www.google.com/');

print $apple_touch_icon_url;
