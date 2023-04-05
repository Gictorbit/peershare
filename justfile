
build: clean
     go build -o ./bin/peershare -ldflags="-s -w" github.com/gictorbit/peershare/cmd
run: build

clean:
     @[ -d "./bin" ] && rm -r ./bin && echo "bin directory cleaned" || true