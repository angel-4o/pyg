from golang:1.22
workdir /app
copy go.mod go.sum ./
run go mod download
copy api/ /app/api
copy cmd/ /app/cmd
run go build -o pyg ./cmd
expose 8080
cmd ["/app/pyg"]