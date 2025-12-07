# **Monolith Backend API — REST + GraphQL (Golang)**

A production-ready backend monolith that provides both **REST API** and **GraphQL API** in a single service.
Built with **Golang**, **Gin/Fiber**, **GQLGen**, **GORM**, and **PostgreSQL**, following **Clean Architecture** and the **Repository Pattern**.

This project is designed as a professional showcase to demonstrate backend engineering skills:
authentication (JWT), role-based access control (RBAC), pagination, filtering, audit logs, and reusable business logic shared between REST and GraphQL.

---

## **Features**

### **Authentication & Authorization**

* Register, login, refresh token
* JWT access & refresh tokens
* Role-Based Access Control (RBAC) — `admin`, `user`
* Context-based auth for REST & GraphQL

### **Product & User Management**

* CRUD operations
* Pagination, filtering, search
* Audit logging (`created_by`, `updated_by`)

### **Dual API (REST + GraphQL)**

* REST: traditional endpoints using Gin/Fiber
* GraphQL: schema-first API using GQLGen
* Both APIs share the same service and repository layer

### **Clean Architecture**

* Delivery layer (REST/GraphQL)
* Service layer (business logic)
* Repository layer (PostgreSQL via GORM)
* Config, middleware, DTO, modular structure

### **Database**

* PostgreSQL (with migrations)
* Soft delete (optional)
* GORM models + repository abstraction

---

## **Tech Stack**

| Layer            | Technology                              |
| ---------------- | --------------------------------------- |
| Web Framework    | Gin / Fiber                             |
| GraphQL          | GQLGen                                  |
| ORM              | GORM                                    |
| Database         | PostgreSQL                              |
| Auth             | JWT                                     |
| Architecture     | Clean Architecture + Repository Pattern |
| Containerization | Docker & Docker Compose                 |

---


## **Running the Project**

### **1. Clone the repository**

```
git clone https://github.com/yourusername/yourproject.git
cd yourproject
```

### **2. Copy environment file**

```
cp .env.example .env
```

### **3. Start PostgreSQL & Backend**

```
docker-compose up --build
```

### **4. Access API**

* REST: `http://localhost:8080/api/v1/...`
* GraphQL Playground: `http://localhost:8080/graphql`

---

## **Environment Variables**

```
APP_PORT=8080
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=secret
DB_NAME=monolith_db
JWT_SECRET=your_jwt_secret
JWT_EXPIRES_IN=15m
REFRESH_EXPIRES_IN=7d
```

---

## **Why This Project Is a Great Backend Showcase**

* Demonstrates mastery of REST + GraphQL in a single monolithic system
* Shows clean architecture, service abstraction, and repository pattern
* Covers authentication, authorization, RBAC, pagination, filtering, audit logs
* Uses GORM, PostgreSQL, JWT—common in real-world backend systems
* Highly maintainable and testable structure
* Ideal to show in interviews and GitHub portfolio
