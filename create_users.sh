#!/bin/sh

BODY='{"name": "Alice Smith", "email": "alice@example.com", "password": "pa55word"}'
curl -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Bob Jones", "email": "bob@example.com", "password": "pa55word"}'
curl -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Carol Smith", "email": "carol@example.com", "password": "pa55word"}'
curl -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Dave Smith", "email": "dave@example.com", "password": "pa55word"}'
curl -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Edith Smith", "email": "edith@example.com", "password": "pa55word"}'
curl -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Faith Smith", "email": "faith@example.com", "password": "pa55word"}'
curl -d "$BODY" localhost:4000/v1/users
