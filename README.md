# Evil Whale
Execute commands inside of docker containers remotely.

## What is Evil Whale?
Evil Whale (or Evil Docker) is an open source red-teaming and pentesting software which attempts to
run commands remotely on a running docker container which is exposed over the internet. Docker uses
a REST API which allows users to modify, create and run containers and more. If this API is exposed
over the internet exploiters and access brokers can exploit the endpoint to gain unauthorized access
to the container and potentially the entire system.
