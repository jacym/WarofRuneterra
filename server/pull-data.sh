#!/bin/sh

dl="https://dd.b.pvp.net/datadragon-set1-lite-en_us.zip"
filename=`basename "$dl"`

if [ ! -f "data/$filename" ]; then
  echo "set1 lite not found, downloading..."
  wget -q -P "$PWD/data" $dl 
fi

if [ -f "data/$filename" ] && [ ! -f "set1-en_us.json" ]; then
  unzip -l "data/$filename" | grep "set1"
  unzip -p "data/$filename" "en_us/data/set1-en_us.json" > "data/set1-en_us.json"
fi

