#!/bin/bash

apt update -y -qq
apt install -y targetcli-fb dbus kmod

# enable d-bus daemon
mkdir /run/dbus
dbus-daemon --system