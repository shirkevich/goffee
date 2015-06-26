#!/bin/bash
exec /go/bin/goffee -phones="$PHONES" -twiliosid="$TWILIOSID" -twiliotoken="$TWILIOTOKEN" -twiliofromnumber="$TWILIOFROMNUMBER" -email="$EMAIL" -slack="$SLACKURL" -clientid="$CLIENTID" -secret="$SECRET" -bind :80 -mandrill="$MANDRILLKEY" -redisaddress="$REDIS_MASTER_SERVICE_HOST:$REDIS_MASTER_SERVICE_PORT" -mysql="$MYSQL" -sessionsecret="$SESSIONSECRET"
