# 🔐 SSH CA Certificate Manager

A **full-stack SSH Certificate Authority manager** with a **Go (Gin) backend** and **Next.js frontend**.  
It provides **secure certificate issuance and access management** with **RBAC, MFA, elliptic curve JWTs, and Vault-based key signing**.  
Built for **performance (Unix sockets)** and **scalability**, with a modern frontend powered by **Next.js + Tailwind + shadcn/ui**.

---

## ✨ Key Features

- 🛡 **RBAC + MFA** — Role-based access control with 6-digit authenticator codes.  
- 🔑 **Elliptic Curve JWTs** — Stateless authentication with secure signing.  
- 📦 **Postgres-backed** — Persistent storage for users, roles, and policies.  
- 🔒 **HashiCorp Vault** — Secure key storage and certificate signing.  
- ⚡ **Unix Sockets** — Faster & safer than TCP for local backend communication.  
- 🎨 **Modern UI** — Built with Next.js, TailwindCSS, and shadcn for reusable components.  
- 🔄 **API Proxying** — Next.js API routes map seamlessly to the Go server.  

---

## 🏗️ Architecture

```text
Frontend (Next.js)                     Backend (Go + Gin)
-----------------                      ------------------------------
- Next.js + Tailwind + shadcn/ui       - Gin REST API (RBAC, MFA)
- API routes proxy to Go API           - JWT (ECC signed) auth
- Reusable UI components               - PostgreSQL (users/roles)
- Secure session handling              - Vault (keys & secrets)
                                       - Unix Sockets for IPC
```
## 🚀 Tech Stack

**Frontend**
- Next.js
- TailwindCSS
- shadcn/ui
- TypeScript

**Backend**
- Go (Gin)
- PostgreSQL
- HashiCorp Vault
- Unix Sockets
- Elliptic Curve Cryptography & ChaCha20-Poly1305 (JWT)

⚙️ Quick Start
1. Clone the repo
```
git clone https://github.com/Dhruvpatel-10/Signee.git
cd Signee
```
2. Setup backend
```
cd ca-api
export USE_UNIX_SOCKET=true
export DATABASE_URL=postgres://user:pass@localhost:5432/sshca
export VAULT_ADDR=http://127.0.0.1:8200
go run main.go
```
3. Setup frontend
```
cd frontend
pnpm install
pnpm run dev
```

🔧 Example Code
Backend: Go (Unix Socket Setup)
```Go
const socketPath = "/tmp/go.sock"

if useSocket := os.Getenv("USE_UNIX_SOCKET"); useSocket == "true" {
    server, listener, err := serveUnixSocket(router, socketPath)
    if err != nil {
        log.Fatalf("unix socket setup error: %v", err)
    }
    server.Serve(listener)
}
```
Frontend: Next.js API Proxy
```ts
const useUnixSocket = process.env.USE_UNIX_SOCKET === "true";
const options: http.RequestOptions = useUnixSocket
  ? { socketPath: "/tmp/go.sock", path, method: req.method, headers: cleanHeaders }
  : { hostname: "localhost", port: process.env.GO_PORT || 8080, path, method: req.method, headers: cleanHeaders };
```

## 🔒 Security Highlights
- Encrypted JWT payloads (user_id, roles) signed with ECC
- Vault-backed secret management and certificate signing
- RBAC + MFA for secure authentication
- Unix sockets reduce attack surface and improve performance

## 🌟 Why This Project?
SSH key management is often manual, insecure, and hard to scale.
This project solves that by providing:

✅ Automated certificate issuance
✅ Centralized role-based access control
✅ Secure cryptographic guarantees
✅ Production-ready frontend for enterprise teams
