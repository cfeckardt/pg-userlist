#!/bin/sh

if [[ -z $OUTPUT_PATH ]];then
    OUTPUT_PATH=/output/userlist.txt
fi

/app -s ' ' -q -l $PG_NOTIFIER -f $OUTPUT_PATH user=$PG_USER password=$PG_PASSWORD host=$PG_HOST port=$PG_PORT dbname=postgres sslmode=disable
