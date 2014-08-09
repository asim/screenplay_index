#!/bin/bash

curl -s http://johnaugust.com/library | grep '\.pdf' | sed -e 's@</a>@\n@g' -e 's@.*href="\([^"]\+\)">@url=\1\&title=@g' | grep title= | sort | uniq
