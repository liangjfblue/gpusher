#!/bin/bash


build-gateway(){
    #build gateway
    echo "build gateway"
    docker build -t gpusher-gateway -f gateway/Dockerfile .
    echo "build gateway ok"
    echo "====================================================="
}

build-logic(){
    #build logic
    echo "build logic"
    docker build -t gpusher-logic -f logic/Dockerfile .
    echo "build logic ok"
    echo "====================================================="
}

build-message(){
    #build message
    echo "build message"
    docker build -t gpusher-message -f message/Dockerfile .
    echo "build message ok"
    echo "====================================================="
}

build-web(){
    #build web
    echo "build web"
    docker build -t gpusher-web -f web/Dockerfile .
    echo "build web ok"
    echo "====================================================="
}


build-all(){
    #build gateway
    build-gateway

    #build logic
    build-logic

    #build message
    build-message

    #build web
    build-web
}

case $1 in
    gateway)
        #build gateway
        build-gateway
        ;;
    logic)
        #build logic
        build-logic
        ;;
    message)
        #build message
        build-message
        ;;
    web)
        #build web
        build-web
        ;;
    all)
        build-all
        ;;
     *)
        echo "Usage: $0 [gateway|logic|message|web|all]"
        exit 1
        ;;
esac