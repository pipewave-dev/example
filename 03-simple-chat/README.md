# How to run example
- Run backend. Ensure port 8080, 29105 and 29106 is idle
```sh
cd ./backend
docker compose up -d
go run .
```
- Run frontend.
```sh
cd ./frontend
pnpm dev
```