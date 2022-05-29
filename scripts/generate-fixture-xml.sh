#!/bin/bash

if [ -z "$1" ]; then
  echo "The script must be given exactly one argument!"
  exit 1
fi

touch ./fixtures/$1-res.xml
touch ./fixtures/$1-req.xml
