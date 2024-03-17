#!/bin/bash

PASS="root"
USER="root"
SQL="./short_url.sql"

out=$(mysql -u$USER -p$PASS < $SQL)

echo $out > /dev/null
