exe := "./bin/pact"

default: build


build:
    go build -o bin/pact

watch-logs:
    tail -f {{cache_directory()}}/pact/logs/demon.log

clean-cache:
    go clean -cache -modcache -i -r

ls-cache:
    ls -lh {{cache_directory()}}/pact

dstart: build
    {{exe}} demon start
    just watch-logs

dstop: build
    {{exe}} demon stop

dstatus: build
    {{exe}} demon status

example *args: build
    EXAMPLE=1 {{exe}} {{args}}

# was used for bubbletea tests
# update-golden-file:
#     go test -v ./... -update
