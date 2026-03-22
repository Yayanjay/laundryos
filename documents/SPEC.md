# LaundryOS MVP - Technical Specification

## 1. Project Overview

| Field             | Value                                                         |
| ----------------- | ------------------------------------------------------------- |
| **Project Name**  | LaundryOS                                                     |
| **Type**          | SaaS Laundry Management System                                |
| **Target Market** | Pemilik usaha laundry (kecil-menengah), Jakarta & Jawa Tengah |
| **Tech Stack**    | React (Next.js 14) + Go (Gin) + PostgreSQL                    |
| **Deployment**    | Cloud SaaS only (AWS Jakarta region)                          |
| **Target Launch** | < 1 bulan (Solo Dev)                                          |

---

## 2. Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        CLIENTS                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐         │
│  │   Browser   │  │   Mobile   │  │   Tablet    │         │
│  │  (Desktop)  │  │   (Future) │  │   (Future)  │         │
│  └──────┬──────┘  └─────────────┘  └─────────────┘         │
└─────────┼───────────────────────────────────────────────────┘
          │ HTTPS
          ▼
┌─────────────────────────────────────────────────────────────┐
│                    FRONTEND (Next.js 14)                    │
│  App Router │ Shadcn/UI │ React Query │ Zod │ Tremor       │
└─────────────────────────┬───────────────────────────────────┘
                          │ REST API (JSON)
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    BACKEND (Go + Gin)                       │
│  Handlers → Services → Repository → DB (PostgreSQL)         │
│       ↑                                    (golang-migrate)  │
│  Middleware: Auth, CORS, Logging, Tenant Context            │
└─────────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────────┐
│                    EXTERNAL SERVICES                        │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   Midtrans   │  │  Sendgrid    │  │   Thermal    │      │
│  │ (QRIS/GoPay) │  │   (Email)    │  │   Printer    │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Multi-Tenancy Strategy

**Approach**: Shared Database dengan `tenant_id` column

```
┌─────────────────────────────────────────────────────┐
│                  tenants table                        │
│  id | name | subdomain | plan | created_at          │
└─────────────────────────────────────────────────────┘
                        │
                        │ 1:N
                        ▼
┌─────────────────────────────────────────────────────┐
│  All other tables have tenant_id FK                 │
│  users | customers | orders | services | payments   │
└─────────────────────────────────────────────────────┘
```

**Middleware Flow**:

1. JWT token berisi `user_id` + `tenant_id`
2. Setiap request, extract tenant_id dari context
3. Semua repository queries filter by tenant_id

---

## 4. Database Schema

### 4.1 Tenants

```sql
CREATE TABLE tenants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    subdomain VARCHAR(100) UNIQUE,
    phone VARCHAR(20),
    address TEXT,
    logo_url VARCHAR(500),
    plan VARCHAR(50) DEFAULT 'starter',
    midtrans_client_key VARCHAR(255),
    midtrans_server_key VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### 4.2 Users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('owner', 'cashier')),
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(tenant_id, email)
);

CREATE INDEX idx_users_tenant ON users(tenant_id);
CREATE INDEX idx_users_email ON users(email);
```

### 4.3 Customers

```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    address TEXT,
    notes TEXT,
    total_orders INT DEFAULT 0,
    total_spent DECIMAL(12,2) DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_customers_tenant ON customers(tenant_id);
CREATE INDEX idx_customers_phone ON customers(phone);
```

### 4.4 Services

```sql
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    unit VARCHAR(50) NOT NULL CHECK (unit IN ('per_kg', 'per_item', 'per_pcs')),
    category VARCHAR(100),
    is_active BOOLEAN DEFAULT true,
    estimated_hours INT DEFAULT 24,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_services_tenant ON services(tenant_id);
```

### 4.5 Orders

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    customer_id UUID REFERENCES customers(id),
    order_number VARCHAR(50) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'new'
        CHECK (status IN ('new', 'processing', 'completed', 'picked_up', 'cancelled')),
    subtotal DECIMAL(10,2) NOT NULL DEFAULT 0,
    discount DECIMAL(10,2) DEFAULT 0,
    discount_type VARCHAR(20) DEFAULT 'nominal',
    final_amount DECIMAL(10,2) NOT NULL DEFAULT 0,
    weight_kg DECIMAL(8,2),
    notes TEXT,
    pickup_date TIMESTAMP WITH TIME ZONE,
    created_by UUID REFERENCES users(id),
    completed_at TIMESTAMP WITH TIME ZONE,
    picked_up_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_orders_tenant ON orders(tenant_id);
CREATE INDEX idx_orders_customer ON orders(customer_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created ON orders(created_at DESC);
CREATE UNIQUE INDEX idx_orders_number_tenant ON orders(tenant_id, order_number);
```

### 4.6 Order Items

```sql
CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    service_id UUID REFERENCES services(id),
    service_name VARCHAR(255) NOT NULL,
    quantity DECIMAL(10,2) NOT NULL DEFAULT 1,
    unit VARCHAR(50) NOT NULL,
    unit_price DECIMAL(10,2) NOT NULL,
    subtotal DECIMAL(10,2) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_order_items_order ON order_items(order_id);
```

### 4.7 Payments

```sql
CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    order_id UUID NOT NULL REFERENCES orders(id),
    amount DECIMAL(10,2) NOT NULL,
    method VARCHAR(50) NOT NULL CHECK (method IN ('cash', 'transfer', 'qris')),
    type VARCHAR(50) DEFAULT 'full',
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed', 'refunded')),
    midtrans_order_id VARCHAR(255),
    midtrans_transaction_id VARCHAR(255),
    midtrans_payment_type VARCHAR(100),
    paid_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_payments_tenant ON payments(tenant_id);
CREATE INDEX idx_payments_order ON payments(order_id);
```

### 4.8 Refresh Tokens

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_token ON refresh_tokens(token_hash);
```

---

## 5. API Endpoints

### 5.1 Authentication

| Method | Endpoint                | Description                          |
| ------ | ----------------------- | ------------------------------------ |
| POST   | `/api/v1/auth/register` | Register tenant + owner              |
| POST   | `/api/v1/auth/login`    | Login, return access + refresh token |
| POST   | `/api/v1/auth/refresh`  | Refresh access token                 |
| POST   | `/api/v1/auth/logout`   | Revoke refresh token                 |
| GET    | `/api/v1/auth/me`       | Get current user                     |

### 5.2 Customers

| Method | Endpoint                   | Description                |
| ------ | -------------------------- | -------------------------- |
| GET    | `/api/v1/customers`        | List customers (paginated) |
| GET    | `/api/v1/customers/:id`    | Get customer detail        |
| POST   | `/api/v1/customers`        | Create customer            |
| PUT    | `/api/v1/customers/:id`    | Update customer            |
| DELETE | `/api/v1/customers/:id`    | Delete customer            |
| GET    | `/api/v1/customers/search` | Search by phone/name       |

### 5.3 Services

| Method | Endpoint               | Description        |
| ------ | ---------------------- | ------------------ |
| GET    | `/api/v1/services`     | List all services  |
| GET    | `/api/v1/services/:id` | Get service detail |
| POST   | `/api/v1/services`     | Create service     |
| PUT    | `/api/v1/services/:id` | Update service     |
| DELETE | `/api/v1/services/:id` | Delete service     |

### 5.4 Orders

| Method | Endpoint                     | Description                |
| ------ | ---------------------------- | -------------------------- |
| GET    | `/api/v1/orders`             | List orders (with filters) |
| GET    | `/api/v1/orders/:id`         | Get order detail           |
| POST   | `/api/v1/orders`             | Create new order           |
| PUT    | `/api/v1/orders/:id`         | Update order               |
| PATCH  | `/api/v1/orders/:id/status`  | Update status              |
| DELETE | `/api/v1/orders/:id`         | Cancel/delete order        |
| GET    | `/api/v1/orders/:id/receipt` | Get receipt data           |

### 5.5 Payments

| Method | Endpoint                             | Description                |
| ------ | ------------------------------------ | -------------------------- |
| POST   | `/api/v1/payments/midtrans/snap`     | Create Midtrans snap token |
| POST   | `/api/v1/payments/midtrans/callback` | Midtrans webhook           |
| POST   | `/api/v1/payments/cash`              | Record cash payment        |

### 5.6 Analytics

| Method | Endpoint                      | Description        |
| ------ | ----------------------------- | ------------------ |
| GET    | `/api/v1/analytics/dashboard` | Dashboard summary  |
| GET    | `/api/v1/analytics/sales`     | Sales by period    |
| GET    | `/api/v1/analytics/orders`    | Order statistics   |
| GET    | `/api/v1/analytics/customers` | Customer analytics |

### 5.7 Settings

| Method | Endpoint            | Description           |
| ------ | ------------------- | --------------------- |
| GET    | `/api/v1/settings`  | Get tenant settings   |
| PUT    | `/api/v1/settings`  | Update settings       |
| GET    | `/api/v1/users`     | List users            |
| POST   | `/api/v1/users`     | Create user (cashier) |
| PUT    | `/api/v1/users/:id` | Update user           |
| DELETE | `/api/v1/users/:id` | Delete user           |

---

## 6. API Response Format

### 6.1 Success Response - With Pagination

```json
{
  "timestamp": "2026-01-31 23:02:53",
  "trace_id": "697e27ad00000000518567838baabee1",
  "response_key": "SUCCESS",
  "message": {
    "title_idn": "",
    "title_eng": "",
    "desc_idn": "",
    "desc_eng": ""
  },
  "data": {
    "items": [ ... ],
    "pagination": {
      "current_page": 1,
      "page_size": 10,
      "total_page": 2,
      "total_item": 20
    }
  }
}
```

### 6.2 Success Response - Without Pagination

```json
{
  "timestamp": "2026-01-31 23:02:53",
  "trace_id": "697e27ad00000000518567838baabee1",
  "response_key": "SUCCESS",
  "message": {
    "title_idn": "",
    "title_eng": "",
    "desc_idn": "",
    "desc_eng": ""
  },
  "data": {
    // object directly, no pagination field
  }
}
```

### 6.3 Error Response

```json
{
  "timestamp": "2026-01-31 23:02:53",
  "trace_id": "697e27ad00000000518567838baabee1",
  "response_key": "ERROR",
  "message": {
    "title_idn": "Terjadi Kesalahan",
    "title_eng": "An Error Occurred",
    "desc_idn": "Pesan error dalam Bahasa Indonesia",
    "desc_eng": "Error message in English"
  },
  "data": null
}
```

### 6.4 Response Keys

| Key       | Usage              |
| --------- | ------------------ |
| `SUCCESS` | Request successful |
| `FAILED`  | General error      |

### 6.5 HTTP Status Codes

| Code | Usage                          |
| ---- | ------------------------------ |
| 200  | Success                        |
| 400  | Bad Request (validation error) |
| 401  | Unauthorized                   |
| 403  | Forbidden                      |
| 404  | Not Found                      |
| 500  | Internal Server Error          |

---

## 7. Project Structure

```
laundry-os/
├── frontend/                          # Next.js 14
│   ├── app/
│   │   ├── (auth)/
│   │   │   ├── login/page.tsx
│   │   │   ├── register/page.tsx
│   │   │   └── layout.tsx
│   │   ├── (dashboard)/
│   │   │   ├── layout.tsx
│   │   │   ├── page.tsx              # Dashboard
│   │   │   ├── orders/
│   │   │   │   ├── page.tsx          # Order list
│   │   │   │   ├── new/page.tsx      # Create order
│   │   │   │   └── [id]/page.tsx     # Order detail
│   │   │   ├── customers/
│   │   │   │   ├── page.tsx          # Customer list
│   │   │   │   └── [id]/page.tsx     # Customer detail
│   │   │   ├── services/
│   │   │   │   ├── page.tsx          # Service list
│   │   │   │   └── [id]/page.tsx    # Edit service
│   │   │   ├── analytics/
│   │   │   │   └── page.tsx          # Analytics dashboard
│   │   │   └── settings/
│   │   │       ├── page.tsx          # General settings
│   │   │       └── users/page.tsx    # User management
│   │   ├── api/
│   │   │   └── [...proxy]/route.ts   # API proxy
│   │   ├── layout.tsx
│   │   └── globals.css
│   ├── components/
│   │   ├── ui/                       # shadcn/ui
│   │   │   ├── button.tsx
│   │   │   ├── input.tsx
│   │   │   ├── select.tsx
│   │   │   ├── table.tsx
│   │   │   ├── card.tsx
│   │   │   ├── dialog.tsx
│   │   │   ├── sheet.tsx
│   │   │   ├── badge.tsx
│   │   │   ├── dropdown-menu.tsx
│   │   │   ├── label.tsx
│   │   │   ├── textarea.tsx
│   │   │   ├── toast.tsx
│   │   │   └── sonner.tsx
│   │   ├── layout/
│   │   │   ├── sidebar.tsx
│   │   │   ├── header.tsx
│   │   │   └── mobile-nav.tsx
│   │   ├── orders/
│   │   │   ├── order-table.tsx
│   │   │   ├── order-form.tsx
│   │   │   ├── order-status-badge.tsx
│   │   │   └── order-receipt.tsx
│   │   ├── customers/
│   │   │   ├── customer-table.tsx
│   │   │   └── customer-form.tsx
│   │   ├── services/
│   │   │   ├── service-table.tsx
│   │   │   └── service-form.tsx
│   │   ├── analytics/
│   │   │   ├── dashboard-stats.tsx
│   │   │   ├── sales-chart.tsx
│   │   │   └── orders-chart.tsx
│   │   └── shared/
│   │       ├── data-table.tsx
│   │       ├── search-input.tsx
│   │       ├── loading-spinner.tsx
│   │       └── empty-state.tsx
│   ├── lib/
│   │   ├── api/
│   │   │   ├── client.ts             # Axios instance
│   │   │   ├── hooks/
│   │   │   │   ├── use-orders.ts
│   │   │   │   ├── use-customers.ts
│   │   │   │   ├── use-services.ts
│   │   │   │   ├── use-analytics.ts
│   │   │   │   └── use-auth.ts
│   │   │   └── mutations/
│   │   │       ├── use-create-order.ts
│   │   │       ├── use-update-order.ts
│   │   │       └── use-payment.ts
│   │   ├── utils.ts
│   │   ├── constants.ts
│   │   └── types.ts
│   ├── .env.local
│   ├── next.config.js
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   └── package.json
│
├── backend/                           # Go
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── api/
│   │   │   ├── handlers/
│   │   │   │   ├── auth.go
│   │   │   │   ├── customer.go
│   │   │   │   ├── service.go
│   │   │   │   ├── order.go
│   │   │   │   ├── payment.go
│   │   │   │   ├── analytics.go
│   │   │   │   └── settings.go
│   │   │   ├── middleware/
│   │   │   │   ├── auth.go
│   │   │   │   ├── tenant.go
│   │   │   │   ├── cors.go
│   │   │   │   └── logger.go
│   │   │   └── router.go
│   │   ├── domain/
│   │   │   ├── tenant.go
│   │   │   ├── user.go
│   │   │   ├── customer.go
│   │   │   ├── service.go
│   │   │   ├── order.go
│   │   │   └── payment.go
│   │   ├── repository/
│   │   │   ├── postgres/
│   │   │   │   ├── tenant.go
│   │   │   │   ├── user.go
│   │   │   │   ├── customer.go
│   │   │   │   ├── service.go
│   │   │   │   ├── order.go
│   │   │   │   └── payment.go
│   │   │   └── interfaces.go
│   │   └── service/
│   │       ├── auth.go
│   │       ├── customer.go
│   │       ├── service.go
│   │       ├── order.go
│   │       ├── payment.go
│   │       └── analytics.go
│   ├── pkg/
│   │   ├── apiresponse/
│   │   │   └── response.go
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── database/
│   │   │   └── postgres.go
│   │   ├── jwt/
│   │   │   └── jwt.go
│   │   ├── midtrans/
│   │   │   └── client.go
│   │   └── hasher/
│   │       └── bcrypt.go
│   ├── migrations/
│   │   ├── 000001_init.up.sql
│   │   └── 000001_init.down.sql
│   ├── .env
│   ├── go.mod
│   └── go.sum
│
├── documents/
│   └── SPEC.md
│
├── .gitignore
├── README.md
└── docker-compose.yml
```

---

## 8. Development Roadmap

### Week 1: Foundation

| Day | Task                                | Deliverable           |
| --- | ----------------------------------- | --------------------- |
| 1   | Project setup (Next.js, Go, Docker) | Dev environment ready |
| 2   | Database design + migrations        | 000001_init.sql       |
| 3   | Auth system - register/login        | JWT + Refresh token   |
| 4   | Auth system - middleware            | Tenant context        |
| 5   | API structure + response format     | Base router           |

### Week 2: Core CRUD

| Day   | Task                | Deliverable         |
| ----- | ------------------- | ------------------- |
| 6-7   | Services CRUD + UI  | Service management  |
| 8-9   | Customers CRUD + UI | Customer management |
| 10-11 | Orders CRUD + UI    | Order management    |
| 12    | Order status flow   | Status transitions  |

### Week 3: Payments & Integration

| Day   | Task                      | Deliverable           |
| ----- | ------------------------- | --------------------- |
| 13-14 | Cash payment + receipt    | Payment recording     |
| 15-16 | Midtrans QRIS integration | Snap token + callback |
| 17    | Receipt generation        | Print-friendly view   |

### Week 4: Analytics & Polish

| Day   | Task                   | Deliverable      |
| ----- | ---------------------- | ---------------- |
| 18-19 | Analytics dashboard    | Tremor charts    |
| 20-21 | Settings + User mgmt   | RBAC             |
| 22-23 | UI polish + responsive | Mobile-friendly  |
| 24-25 | Bug fixes + testing    | QA               |
| 26-28 | Deployment + docs      | Production-ready |

---

## 9. User Personas

### Persona 1: Owner

| Attribute       | Description                                                       |
| --------------- | ----------------------------------------------------------------- |
| **Role**        | Pemilik laundry, oversight penuh                                  |
| **Goals**       | Monitor bisnis, analisis profit, kontrol biaya                    |
| **Pain Points** | Tidak bisa lihat laporan real-time, tidak tahu profit per layanan |
| **Needs**       | Dashboard analytics, multi-user access, reports                   |

### Persona 2: Kasir

| Attribute       | Description                                               |
| --------------- | --------------------------------------------------------- |
| **Role**        | Operator frontline, high volume transactions              |
| **Goals**       | Input order cepat, tidak salah hitung, receipt clear      |
| **Pain Points** | Excel prone to error, lost receipts, manual calculation   |
| **Needs**       | Fast UI, auto-calculation, print receipt, status tracking |

### Persona 3: Pelanggan

| Attribute       | Description                                     |
| --------------- | ----------------------------------------------- |
| **Role**        | End customer, ingin tahu status laundry         |
| **Goals**       | Tahu kapan selesai, tidak lupa ambil            |
| **Pain Points** | Tidak ada notifikasi, harus telp/tanya langsung |
| **Needs**       | Simple status check, pickup reminder            |

---

## 10. Monetization Strategy

### Pricing Tiers

| Plan           | Harga/Bulan | Target            | Features                                             |
| -------------- | ----------- | ----------------- | ---------------------------------------------------- |
| **Starter**    | Rp 99.000   | Laundry kaki lima | 1 outlet, 2 users, 300 orders/mo, basic reports      |
| **Business**   | Rp 249.000  | Laundry menengah  | 1 outlet, 5 users, unlimited orders, analytics, QRIS |
| **Enterprise** | Rp 499.000  | Franchise         | 3 outlets, unlimited users, multi-outlet reports     |

### Additional Revenue

- Add-on: Extra outlet (Rp 150.000/outlet/mo)
- Add-on: SMS notification (Rp 0.05/SMS)
- Transaction fee: 0% for Business+, 1% for Starter

---

## 11. Configuration

### Backend Environment Variables

```bash
# Server
PORT=8080
ENV=development

# Database
DATABASE_URL=postgres://user:pass@host:5432/laundryos?sslmode=disable

# JWT
JWT_SECRET=<256-bit-secret>
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h  # 7 days

# Midtrans
MIDTRANS_IS_PRODUCTION=false
MIDTRANS_SERVER_KEY=<key>
MIDTRANS_CLIENT_KEY=<key>

# CORS
CORS_ORIGIN=http://localhost:3000
```

### Frontend Environment Variables

```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
NEXT_PUBLIC_APP_NAME=LaundryOS
```

---

## 12. Acceptance Criteria

### Must Have (MVP)

- [ ] Registration + login dengan JWT
- [ ] CRUD Services dengan pricing
- [ ] CRUD Customers dengan search
- [ ] Create Order dengan multiple items
- [ ] Order status flow (new → processing → completed → picked_up)
- [ ] Cash payment recording
- [ ] Receipt generation (printable)
- [ ] Dashboard dengan daily sales
- [ ] Basic analytics (sales trend, order count)
- [ ] Responsive UI untuk desktop/tablet

### Should Have

- [ ] Midtrans QRIS integration
- [ ] Discount (nominal/percent)
- [ ] Order number auto-generation
- [ ] Customer history view

### Nice to Have (Post-MVP)

- [ ] WhatsApp notification
- [ ] Multi-outlet support
- [ ] Inventory tracking
- [ ] Mobile app
- [ ] Customer portal

---

## 13. Tech Dependencies

### Frontend

```json
{
  "next": "14.x",
  "react": "18.x",
  "@tanstack/react-query": "5.x",
  "react-hook-form": "7.x",
  "zod": "3.x",
  "@hookform/resolvers": "3.x",
  "axios": "1.x",
  "sonner": "1.x",
  "@tremor/react": "3.x",
  "date-fns": "3.x",
  "lucide-react": "latest",
  "tailwindcss": "3.x",
  "class-variance-authority": "latest",
  "clsx": "latest",
  "tailwind-merge": "latest",
  "@radix-ui/react-*": "latest"
}
```

### Backend

```go
// go.mod
require (
    github.com/gin-gonic/gin
    github.com/lib/pq
    github.com/golang-jwt/jwt/v5
    github.com/google/uuid
    golang.org/x/crypto
    github.com/midtrans/midtrans-go
    github.com/jmoiron/sqlx
    github.com/golang-migrate/migrate/v4
)
```

---

## 14. Security Considerations

1. **Password**: bcrypt with cost 12
2. **JWT**: HS256, short-lived access token (15m)
3. **Refresh Token**: Stored as hash in DB, revocable
4. **Tenant Isolation**: All queries filtered by tenant_id from JWT
5. **Input Validation**: Zod (frontend) + manual (backend)
6. **SQL Injection**: Use parameterized queries only
7. **CORS**: Strict origin checking
8. **HTTPS**: Mandatory in production

---

_Document Version: 1.0_
_Last Updated: 2026-03-22_
_Author: Product Specification_
