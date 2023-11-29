#!/bin/sh

echo $1

BODY='{"title":"Moana","year":2016,"runtime":"107 mins", "genres":["animation","adventure"]}'
curl -d "$BODY" localhost:4000/v1/movies -H "Authorization: Bearer $1"

BODY='{"title":"Black Panther","year":2018,"runtime":"134 mins","genres":["action","adventure"]}'
curl -d "$BODY" localhost:4000/v1/movies -H "Authorization: Bearer $1"

BODY='{"title":"Deadpool","year":2016, "runtime":"108 mins","genres":["action","comedy"]}'
curl -d "$BODY" localhost:4000/v1/movies -H "Authorization: Bearer $1"

BODY='{"title":"The Breakfast Club","year":1986, "runtime":"96 mins","genres":["drama"]}'
curl -d "$BODY" localhost:4000/v1/movies -H "Authorization: Bearer $1"
