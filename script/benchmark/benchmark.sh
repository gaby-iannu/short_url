#!/bin/bash
JSON_FILE="./urls.json"

i=0
while read -r LINE
do
	echo $LINE > "$i.json"
	`ab -p $i.json -T application/json -c 5 -n 15 127.0.0.1:8080/tiny > $i.log &`

	i=$(($i+1))

	if [ $((i % 30)) -eq 0 ]; then
		sleep 1
	fi

done < $JSON_FILE

while ![ -f $i.log ]
do
	sleep 1
done

LOG="benchmark.log"
for j in $(seq 0 $i)
do
	cat $j.json >> $LOG 
	`rm $j.json $j.log`
done
