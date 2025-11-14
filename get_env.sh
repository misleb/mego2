#!/bin/bash

export GOOGLE_CLIENT_ID=$(heroku config:get GOOGLE_CLIENT_ID)
export GOOGLE_CLIENT_SECRET=$(heroku config:get GOOGLE_CLIENT_SECRET)
export DATABASE_URL=$(heroku config:get DATABASE_URL)
export BASE_URI=http://localhost:8080
export TEST_DATABASE_URL=postgresql://mego2:mego2_dev@localhost:5432/mego2_test?sslmode=disable
