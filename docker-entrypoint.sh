wait-for "${DATABASE_HOST}:${DATABASE_PORT}"  -- "$@"

#watch for .go file and invoke go build if the files changes changes
#CompileDaemon --build="go build -o main main.go" --command=./app/main