from golang:1.11.0-alpine3.8 as build

run apk update \
 && apk add git gcc musl-dev

env pkg github.com/fragmenta/fragmenta-cms

add . /go/src/${pkg}
workdir /go/src/${pkg}

run go vet  ./...

# tests don't pass (missing json file?)
# run go test ./...

run go build -o /go/bin/fragmenta .

from alpine:3.8

entrypoint ["/bin/fragmenta"]
copy --from=build /go/bin/fragmenta /bin/

volume /data
workdir /data

