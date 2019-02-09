#!/bin/bash
rm ./esnextnews ./esnextnews.zip
GOOS=linux go build esnextnews.go
zip esnextnews.zip ./esnextnews