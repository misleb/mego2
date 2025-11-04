#!/bin/bash

export GOOGLE_CLIENT_ID=$(heroku config:get GOOGLE_CLIENT_ID)
export GOOGLE_CLIENT_SECRET=$(heroku config:get GOOGLE_CLIENT_SECRET)
export DATABASE_URL=$(heroku config:get DATABASE_URL)
export BASE_URI=http://localhost:8080