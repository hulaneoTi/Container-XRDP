#!/bin/sh

if [ -r /etc/default/locale ]; then
  . /etc/default/locale
  export LANG LANGUAGE
fi

# Default
#. /etc/X11/Xsession

# XFCE
startxfce4

#openbox & setxkbmap -layout "br" &
#sleep 1
#/opt/firefox/firefox -width 1920 -height 1080 & wait $!
#pkill -TERM -u $USER firefox
#pkill -KILL -u $USER