exe := "./bin/pact"

default: build

watch-logs:
    tail -f $XDG_CACHE_HOME/pact/logs/demon.log

update-golden-file:
    go test -v ./... -update

build:
    go build -o bin/pact

clean-cache:
    go clean -cache -modcache -i -r

ls-cache:
    ls -lh $XDG_CACHE_HOME/pact

dstart: build
    {{exe}} demon start
    just watch-logs

dstop: build
    {{exe}} demon stop

dstatus: build
    {{exe}} demon status
