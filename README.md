## Morgue

To store and request response bodies.

It uses the URI as the storage key.

Store response, with JSON as content type

    curl -v -XPUT -H 'Content-Type:application/json' -d "{'temperature': 30}" 0.0.0.0:8080/aa/bb/cc

Request

    curl http://0.0.0.0:8080/aa/bb/cc
    {'temperature': 30}

Store response with two seconds of sleep time

    curl -vv -XPUT -H 'morgue-set-sleep-time:2000' -d "{'temperature': 30}" 0.0.0.0:8080/aa/bb/cc

### Install

    go get -u github.com/juanpabloaj/morgue/cmd/morgue

Change the default port with the environment variable `PORT`

    export PORT=8080

### Run docker image

    docker run --rm --name morgue -p 0.0.0.0:8080:8080 juanpabloaj/morgue:latest
