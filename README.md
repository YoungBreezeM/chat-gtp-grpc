# Chat GTP Grpc

## Overview
Chat-Gtp-grpc program is a chat gtp free api with grpc protocol. Supportting a simple function to use. 

## Build

1. Build a binary file

    ```shell
    make build
    ```
2. Build done,you will has a server binary

## RUN

1. Normal to start

    ```shell
     server
    ```
    Change server port.

    ```shell
     SERVER_PORT = "8080" server
    ```
    If you want to chat gtp through proxy.

    ```shell
     PROXY_URL = "socket5:127.0.0.1" server
    ```