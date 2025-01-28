#/bin/bash
set -e;

function onExit {
    if [ "$?" != "0" ]; then
        echo "Tests failed";
        # build failed, don't deploy
        exit 1;
    else
        echo "Tests passed";
        # deploy build
    fi
}

trap onExit EXIT;

docker run -t postman/newman:alpine https://www.getpostman.com/collections/8a0c9bc08f062d12dcda --suppress-exit-code;
