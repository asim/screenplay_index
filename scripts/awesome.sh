#!/bin/bash

curl -s  http://www.awesomefilm.com/ | egrep '\.pdf' | sed 's@</a>@\n@g' | egrep '\.pdf' | sed -e 's@^.*\(href\|HREF\)="@url=http://www.awesomefilm.com/@g' -e 's/">/\&title=/g'  |  grep ^url
