#!/bin/bash

for j in *; do m=`cat -- $j| sed 's@"@\\\"@g'`; echo -e "{\"doc\" : { \"meta\" : \"$m\" } }" > t; curl -T t -XPOST "http://127.0.0.1:9200/scripts/script/$j/_update" ; done
