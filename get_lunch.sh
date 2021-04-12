#!/bin/bash
# https://www.tu.no/?module=TekComics&service=image&id=lunch&key=2020-04-14
url="https://www.tu.no/"
date="$1"
if [ -z "$date" ]
then
    date=$(date +%F)
fi
key=$date
module="TekComics"
service="image"

id="lunch"
destination="~/lunch/tu-$id-$date.jpg"
# https://www.tu.no/?module=TekComics&service=image&id=lunch&key=2020-05-28
curl --get --url "$url" --data-urlencode "module=$module"  --data-urlencode "service=$service" --data-urlencode "id=$id" --data-urlencode "key=$date" --output "$destination"
cd ~/lunch
find . -name 'tu-$id-*' -size 0 -print0 | xargs -0 rm

