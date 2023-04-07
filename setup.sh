#!/bin/bash

apt update
apt install golang

sourceList=".sources.txt"
allowPath=".allow.txt"
blockPath=".block.txt"

s1="https://adaway.org/hosts.txt"

s2="https://pgl.yoyo.org/adservers/serverlist.php?hostformat=hosts&showintro=0&mimetype=plaintext"

# This link may not work in India
s3="# https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts"

if [ ! -f $allowPath ]; then
    touch $allowPath
fi

if [ ! -f $blockPath ]; then
    touch $blockPath
fi

if [ ! -f $sourceList ]; then
    echo $s1 >> $sourceList
    echo $s2 >> $sourceList
    echo $s3 >> $sourceList
fi

go build .
