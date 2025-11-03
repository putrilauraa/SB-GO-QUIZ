# SB-GO-QUIZ

Mini Project: Quiz Bootcamp Golang
Project RESTful API untuk data buku dan kategori menggunakan Golang dengan framework Gin dan PostgreSQL.

Local Setup
Dependencies: Golang, PostgreSQL.
Packages: Gin , lib/pq , sql-migrate.
Project inititaliation:
go mod init your_module_name
go get github.com/gin-gonic/gin
go get github.com/lib/pq
go install github.com/rubenv/sql-migrate/...@latest
Database Migration:
sql-migrate up -config=dbconfig.yml -env=production

Deployment & Environment
Deployment Platform: Vercel (Go) & Railway (Postgres).

URL Vercel: https://sb-go-quiz.vercel.app/
Environment Variables:
DATABASE_URL: postgresql://postgres:GFDzLLbKVMlygrjfUTPYIglqiCOSQlZo@trolley.proxy.rlwy.net:18131/railway

Authentication (Basic Auth)
Method: Middleware Basic Auth 
Testing Credential:
Username: admin
Password: password123
Include header Authorization: Basic YWRtaW46cGFzc3dvcmQxMjM= on every request to /api/*.

API Endpoints List
Feature: Category
GET /api/categories (Menampilkan seluruh kategori)
POST /api/categories (Menambahkan kategori)
GET /api/categories/:id (Detail kategori)
DELETE /api/categories/:id (Menghapus kategori)
GET /api/categories/:id/books (Buku berdasarkan kategori)
Feature: Book
GET /api/books (Menampilkan seluruh buku)
POST /api/books (Menambahkan buku)
GET /api/books/:id (Detail buku)
DELETE /api/books/:id (Menghapus buku)