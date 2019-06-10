#!/bin/bash

export TOKEN=tst
export PORT=3090

gin --appPort 3090 --port 3097 -i -x .git -x vendor --build ./ run
