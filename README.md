# Event Management System

A REST API-based event management system built with Go (Golang) using the Gin framework. This application provides comprehensive features for creating, managing, and registering events with integrated payment processing via Midtrans.

## 🚀 Key Features

### Authentication & Authorization

- **User Registration** with email verification using OTP
- **Login** with email & password
- **OAuth 2.0** with Google Sign-In
- **JWT-based Authentication** for API security
- **Role-based Access Control** (Admin & User)

### Event Management

- **Event Creation** with complete details (title, description, category, location, quota, price)
- **Location Types** supporting both offline and online events
- **Event Status** (draft, published, cancelled)
- **Session Management** for multi-session events

### Registration & Payment System

- **Event Registration** with automatic ticket generation
- **Midtrans Integration** for payment gateway
- **Multiple Payment Methods** (e-wallet, bank transfer, credit card)
- **Webhook Notifications** for real-time payment status updates
- **Automatic Ticket Code Generation** for each registration

## 🏗️ Architecture

This project follows **Clean Architecture** with clear layer separation:

```
event-management/
├── cmd/app/              # Application entry point
├── entity/               # Domain entities (database models)
├── model/                # DTOs (Data Transfer Objects)
├── internal/
│   ├── handler/rest/     # HTTP handlers (controllers)
│   ├── service/          # Business logic layer
│   └── repository/       # Data access layer
└── pkg/
    ├── bcrypt/           # Password hashing utility
    ├── config/           # Configuration management
    ├── database/         # Database connection & migration
    ├── jwt/              # JWT token management
    ├── mail/             # Email service
    ├── middleware/       # HTTP middlewares
    ├── response/         # Standardized API responses
    └── supabase/         # Supabase integration
```

## 📦 Tech Stack

### Core Framework & Language

- **Go 1.24.1** - Programming language
- **Gin** - HTTP web framework
- **GORM** - ORM library for database operations

### Database

- **MariaDB/MySQL** - Relational database

### Third-Party Integrations

- **Midtrans** - Payment gateway (Snap & Core API)
- **Google OAuth 2.0** - Social authentication
- **Supabase Storage** - File storage (optional)
- **SMTP** - Email service for OTP delivery

### Security & Authentication

- **JWT (golang-jwt/jwt)** - Token-based authentication
- **bcrypt** - Password hashing
- **CORS** - Cross-Origin Resource Sharing

## 🗄️ Database Schema

### Entities

**Users**

- User ID (UUID)
- Role ID (Admin/User)
- Google ID (for OAuth)
- Name, Email, Password
- Profile Picture
- Account Status (active/inactive)

**Events**

- Event ID (UUID)
- User ID (creator)
- Title, Description, Category
- Start Date, End Date
- Location & Location Type (offline/online)
- Quota & Price
- Status (draft/published/cancelled)

**Sessions**

- Session ID (UUID)
- Event ID
- Title, Start Time, End Time

**Registrations**

- Registration ID (UUID)
- Event ID, User ID
- Ticket Code (6 unique characters)
- Status (pending/approved/rejected)

**Payments**

- Order ID (generated)
- Registration ID
- Amount, Status
- Snap URL (Midtrans payment link)
- Payment Type
- Paid At timestamp

**OTP**

- OTP ID (UUID)
- User ID
- Code (6 digits)
- Expiration handling

**Roles**

- Role ID
- Role Name

## 🔧 Setup & Installation

### Prerequisites

- Go 1.24.1 or higher
- MariaDB/MySQL
- Midtrans Account (Sandbox/Production)
- Google OAuth 2.0 Credentials
- SMTP Server for email

### Environment Variables

Create a `.env` file in the root directory with the following configuration:

```env
# Server Configuration
ADDRESS=localhost
PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=event_management

# JWT Configuration
JWT_SECRET=your_jwt_secret_key

# OTP Configuration
EXPIRED_OTP=5  # in minutes

# Email Configuration
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your_email@gmail.com
SMTP_PASSWORD=your_app_password

# Google OAuth
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback

# Midtrans Configuration
MIDTRANS_SERVER_KEY=your_midtrans_server_key
MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_ENVIRONMENT=sandbox  # or production
```

### Installation Steps

1. **Clone the repository**

```bash
git clone <repository-url>
cd event-management
```

2. **Install dependencies**

```bash
go mod download
```

3. **Setup database**

- Create a new database in MariaDB/MySQL
- Update database configuration in `.env`

4. **Run migrations**

```bash
go run cmd/app/main.go
```

Migrations will run automatically when the application starts

5. **Run the application**

```bash
go run cmd/app/main.go
```

The server will run at `http://localhost:8080`

## 📡 API Endpoints

### Authentication

- `POST /api/v1/auth/register` - Register new user
- `PATCH /api/v1/auth/register` - Verify OTP
- `PATCH /api/v1/auth/register/resend` - Resend OTP
- `POST /api/v1/auth/login` - Login with email & password
- `GET /api/v1/auth/google/login` - Initiate Google OAuth
- `GET /api/v1/auth/google/callback` - Google OAuth callback

### Events (Protected)

- `POST /api/v1/events/add-event` - Create new event

### Registrations (Protected)

- `POST /api/v1/registrations/event` - Register for an event

### Notifications

- `POST /api/v1/notification` - Webhook for Midtrans payment notifications

## 🔐 Middleware

- **CORS** - Allows cross-origin requests
- **Authentication** - JWT token validation
- **Authorization** - Role-based access control
- **Timeout** - Request timeout handling

## 💳 Payment Flow

1. User registers for an event
2. System creates registration with "pending" status
3. System generates payment via Midtrans Snap
4. User receives Snap URL for payment
5. User completes payment
6. Midtrans sends webhook notification
7. System updates payment and registration status to "approved"
8. User receives ticket code

## 📧 Email Verification Flow

1. User registers with email & password
2. System generates 6-digit OTP
3. OTP is sent to user's email
4. User verifies with OTP within the specified time (default: 5 minutes)
5. Account status changes to "active"

## 🔒 Security Features

- **Password Hashing** using bcrypt
- **JWT Tokens** for stateless authentication
- **OTP Expiration** for email verification security
- **CORS Protection** to prevent unauthorized access
- **Request Timeout** to prevent DoS attacks
- **Transaction Management** for data consistency

## 🛠️ Development

### Project Structure Explanation

- **entity/** - Database table representations with GORM tags
- **model/** - Request/Response DTOs for API
- **internal/handler/** - HTTP request handlers
- **internal/service/** - Business logic & orchestration
- **internal/repository/** - Database queries & operations
- **pkg/** - Reusable utilities & packages

### Adding New Features

1. Define entity in `entity/`
2. Create request/response models in `model/`
3. Implement repository in `internal/repository/`
4. Implement business logic in `internal/service/`
5. Create handler in `internal/handler/rest/`
6. Register endpoint in `rest.go`
