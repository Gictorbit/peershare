
# format and lint and generate proto file using buf
proto:
    rm -rf api/gen
    @echo "run proto linter..."
    @cd api && buf lint && cd -

    @echo "format proto..."
    @cd api && buf format -w && cd -

    @echo "generate proto..."
    @cd api && buf generate && cd -

build: clean
     go build -o ./bin/peershare -ldflags="-s -w" github.com/Gictorbit/peershare/cmd
run: build

clean:
     @[ -d "./bin" ] && rm -r ./bin && echo "bin directory cleaned" || true