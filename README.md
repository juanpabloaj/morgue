## Morgue

To store and request body responses.

It uses the URI as the storage key.

Store

    curl -vv -XPUT -d "{'temperature': 30}" http://0.0.0.0:8080/aa/bb/cc

Request

    curl http://0.0.0.0:8080/aa/bb/cc
    {'temperature': 30}

Run docker image

    docker run --rm --name morgue -p 0.0.0.0:8080:8080 juanpabloaj/morgue:latest
