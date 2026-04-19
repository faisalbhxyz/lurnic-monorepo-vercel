# Project Handover & Technical Documentation

## 1. Project Overview
**Project Name:** Dashlearn Server (Lurnic)

**Description:** A multi-tenant learning management system (LMS) backend that handles course creation, student enrollments, order processing, and quizzes/assignments. It allows instructors to manage courses and students to enroll and track their progress, while supporting dedicated tenant environments.

**Tech Stack:**
* **Language/Framework:** Go (Golang 1.24) with Gin Web Framework
* **Database:** MySQL
* **ORM & Migrations:** GORM and Goose
* **Auth:** JWT (JSON Web Tokens)
* **Storage/Cloud:** DigitalOcean Spaces / Cloudflare R2 (S3-compatible via AWS SDK for Go)

## 2. Infrastructure & Access
**Repository:** [abukhalidrifat/dashlearn-server](https://github.com/abukhalidrifat/dashlearn-server)

**Hosting:** Needs to be configured (e.g., DigitalOcean App Platform, Render, AWS, Hetzner VPS). It includes a `Dockerfile` and `compose.yaml` for containerized deployments.

**Domain Provider:** To be configured.

**Database Host:** MySQL Database (Local or Managed like AWS RDS, DigitalOcean Managed Databases).

## 3. Environment Variables (.env)
Here is the required format for the environment variables. Keep the actual keys secure.

```env
APP_PORT=5000
GOOSE_MIGRATION_DIR=migrations
GOOSE_DRIVER=mysql
GOOSE_DBSTRING=root:@tcp(127.0.0.1:3306)/dashlearn?charset=utf8mb4&parseTime=True&loc=Local
JWT_SECRET=your_jwt_secret_here

# Cloudflare R2 Configuration (S3-compatible)
R2_PUBLIC_BASE_URL=https://cdn.lurnic.com
R2_ACCOUNT_ID=your_cloudflare_account_id
R2_ACCESS_KEY_ID=your_r2_access_key_id
R2_SECRET_ACCESS_KEY=your_r2_secret_access_key
R2_BUCKET=lurnic
R2_UPLOAD_PREFIX=
R2_REGION=auto
```

## 4. Database Schema
**Architecture:** Relational Database model using MySQL. Highly structured around **Multi-Tenancy**, where almost all tables have a `tenant_id` foreign key referencing the `tenants` table.

**Key Tables:**
* **tenants**: The core table for multi-tenancy.
* **users / roles**: User profiles and RBAC (Role-Based Access Control) with specific permissions. Users are linked to roles.
* **course_details**: Core information of courses. Created by an `author_id` (users).
* **course_chapters / course_lessons**: Course content hierarchy. Includes relations for resources (lesson_resources), quizzes, and assignments.
* **students / instructors**: Specialized profiles linked to courses.
* **enrollment**: Pivot tracking which `student_id` is enrolled in which `course_id`.
* **orders / payment_methods**: Transaction tracking and history. Orders refer to `student_id` and `course_id`.

**Relationships:**
* `users` belong to `roles` (`role_id` -> `roles(id)`).
* `course_details`, `roles`, `categories`, `students`, `orders`, `enrollment`, etc. belong to a `tenant` (`tenant_id` -> `tenants(id)`) ensuring isolated data per tenant.
* `course_chapters` map to `course_details` (`course_id` -> `course_details(id)`).
* `course_lessons` map to `course_chapters` (`chapter_id` -> `course_chapters(id)`).
* `enrollment` and `orders` join `students` and `course_details`.
* `course_quizzes` and `course_assignments` belong to both a `course` and a specific `chapter`.

## 5. Core Functionalities (Logic Flow)
* **Multi-Tenancy & Authorization Flow:**
  API requests are isolated by tenant context. Authentication validates JWT tokens via middleware and extracts user context and role. Roles define RBAC permissions (`roles.permissions`).
* **Course Management Flow:**
  Instructors (`user` module -> `instructor` module) create `course` records (which encompass `chapters`, `lessons`, and `lesson_resources`). Settings and categories are managed separately but mapped to courses.
* **Student Enrollment & Payment Flow:**
  When a student purchases a course, the `order` module processes the transaction (storing `transaction_id`), and upon success, an entry is created in the `enrollment` table, granting the student access to the course content.
* **File Uploads & Storage:**
  Files (like banners or course images) can use S3-compatible endpoints, integrating either DigitalOcean Spaces or Cloudflare R2.

## 6. Local Setup Guide
1. **Clone the repository:**
   ```bash
   git clone https://github.com/abukhalidrifat/dashlearn-server.git
   cd dashlearn-server
   ```
2. **Install Go Modules:**
   ```bash
   go mod download
   # or
   go mod tidy
   ```
3. **Environment Setup:**
   * Copy `.env.example` to `.env` (or create a `.env` file).
   * Update the `GOOSE_DBSTRING` to match your local MySQL server setup. Ensure you create a database named `dashlearn` in your MySQL server.
4. **Run Database Migrations:**
   Install goose if you haven't (`go install github.com/pressly/goose/v3/cmd/goose@latest`).
   ```bash
   goose -dir migrations mysql "root:@tcp(127.0.0.1:3306)/dashlearn?charset=utf8mb4&parseTime=True&loc=Local" up
   ```
5. **Run the Project Locally:**
   ```bash
   go run main.go
   ```
   The backend will be available at `http://localhost:5000`.

*Alternatively, use Docker Compose (full stack: API + dashboard):*
```bash
docker compose -f docker-compose.yaml up -d --build
```
