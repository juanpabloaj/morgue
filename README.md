## Morgue

To store and request body responses.

It uses the URI as the storage key.

Store

    curl -vv -XPUT -d "{'temperature': 30}" http://0.0.0.0:8080/aa/bb/cc

Request

    curl http://0.0.0.0:8080/aa/bb/cc
    {'temperature': 30}
