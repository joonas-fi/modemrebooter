FROM fn61/buildkit-golang:20181204_1302_5eedb86addc826e7

WORKDIR /go/src/github.com/joonas-fi/modemrebooter

CMD bin/build.sh
