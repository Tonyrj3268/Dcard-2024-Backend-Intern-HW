#!/bin/bash

source .env

PGPASSWORD="$POSTGRES_PASSWORD" psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -h "$POSTGRES_HOSTNAME"\
     -f ./pg_sql/set_trigger.sql