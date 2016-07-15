#!/usr/bin/env bash
DIR=`dirname $0`
cqlsh -e "DROP TABLE athena.matches;"
cqlsh -f $DIR/../db/schema/matches.cql
cqlsh -e "DROP TABLE athena.rankings;"
cqlsh -f $DIR/../db/schema/rankings.cql
echo "Done."
