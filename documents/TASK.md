# LaundryOS - Task Board (Living Document)

> Last Updated: 2026-03-22
> Status: `IN_PROGRESS` | `TODO` | `DONE`

---

## Phases Overview

| Phase | Name | Status | Progress |
|-------|------|--------|----------|
| 0 | Project Setup | DONE | 100% |
| 1 | Foundation | TODO | 0% |
| 2 | Core Features | TODO | 0% |
| 3 | Payments | TODO | 0% |
| 4 | Analytics | TODO | 0% |
| 5 | Polish & Deploy | TODO | 0% |

---

## Phase 0: Project Setup

### 0.1 Infrastructure
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 0.1.1 | Setup project directory structure | DONE | HIGH | 1h |
| 0.1.2 | Create docker-compose.yml (PostgreSQL + Redis) | DONE | HIGH | 1h |
| 0.1.3 | Create .env files (backend, frontend) | DONE | HIGH | 30m |
| 0.1.4 | Setup Git repository + .gitignore | DONE | MEDIUM | 15m |

### 0.2 Backend Setup
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 0.2.1 | Initialize Go module | DONE | HIGH | 15m |
| 0.2.2 | Setup standard Go layout (cmd/, internal/, pkg/) | DONE | HIGH | 30m |
| 0.2.3 | Install dependencies (gin, lib/pq, jwt, bcrypt, uuid, migrate) | DONE | HIGH | 30m |
| 0.2.4 | Create config loader | DONE | HIGH | 1h |
| 0.2.5 | Create database connection package | DONE | HIGH | 1h |

### 0.3 Frontend Setup
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 0.3.1 | Initialize Next.js 14 project | DONE | HIGH | 15m |
| 0.3.2 | Setup Tailwind CSS + shadcn/ui | DONE | HIGH | 2h |
| 0.3.3 | Install dependencies (react-query, zod, axios, tremor, date-fns) | DONE | HIGH | 1h |
| 0.3.4 | Setup project structure (app router folders) | DONE | HIGH | 1h |
| 0.3.5 | Create base components (ui/ folder) | DONE | HIGH | 4h |
| 0.3.6 | Create API client (axios instance) | DONE | HIGH | 1h |

### 0.4 Migrations
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 0.4.1 | Create 000001_init.up.sql | DONE | HIGH | 2h |
| 0.4.2 | Create 000001_init.down.sql | DONE | HIGH | 30m |
| 0.4.3 | Setup golang-migrate | DONE | HIGH | 30m |
| 0.4.4 | Test migration up/down | DONE | HIGH | 30m |

---

## Phase 1: Foundation

### 1.1 API Response Package
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 1.1.1 | Create API response struct (pkg/apiresponse) | TODO | HIGH | 1h |
| 1.1.2 | Implement Success response helper | TODO | HIGH | 30m |
| 1.1.3 | Implement Error response helper | TODO | HIGH | 30m |
| 1.1.4 | Implement Pagination helper | TODO | HIGH | 1h |
| 1.1.5 | Implement Trace ID middleware | TODO | HIGH | 1h |

### 1.2 JWT & Auth
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 1.2.1 | Create JWT package (pkg/jwt) | TODO | HIGH | 2h |
| 1.2.2 | Implement token generation (access + refresh) | TODO | HIGH | 2h |
| 1.2.3 | Implement token validation | TODO | HIGH | 1h |
| 1.2.4 | Implement refresh token rotation | TODO | HIGH | 2h |
| 1.2.5 | Create auth middleware | TODO | HIGH | 2h |
| 1.2.6 | Create tenant context middleware | TODO | HIGH | 2h |

### 1.3 Auth Handler & Routes
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 1.3.1 | POST /api/v1/auth/register | TODO | HIGH | 3h |
| 1.3.2 | POST /api/v1/auth/login | TODO | HIGH | 3h |
| 1.3.3 | POST /api/v1/auth/refresh | TODO | HIGH | 2h |
| 1.3.4 | POST /api/v1/auth/logout | TODO | HIGH | 1h |
| 1.3.5 | GET /api/v1/auth/me | TODO | HIGH | 1h |

### 1.4 User Management
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 1.4.1 | User repository (CRUD) | TODO | HIGH | 3h |
| 1.4.2 | User service layer | TODO | HIGH | 2h |
| 1.4.3 | GET /api/v1/users | TODO | HIGH | 1h |
| 1.4.4 | POST /api/v1/users | TODO | HIGH | 2h |
| 1.4.5 | PUT /api/v1/users/:id | TODO | HIGH | 2h |
| 1.4.6 | DELETE /api/v1/users/:id | TODO | HIGH | 1h |

### 1.5 Tenant Settings
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 1.5.1 | Tenant repository | TODO | HIGH | 2h |
| 1.5.2 | GET /api/v1/settings | TODO | HIGH | 1h |
| 1.5.3 | PUT /api/v1/settings | TODO | HIGH | 2h |

---

## Phase 2: Core Features

### 2.1 Services Management
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 2.1.1 | Service domain + repository | TODO | HIGH | 2h |
| 2.1.2 | Service service layer | TODO | HIGH | 2h |
| 2.1.3 | GET /api/v1/services | TODO | HIGH | 1h |
| 2.1.4 | GET /api/v1/services/:id | TODO | HIGH | 1h |
| 2.1.5 | POST /api/v1/services | TODO | HIGH | 2h |
| 2.1.6 | PUT /api/v1/services/:id | TODO | HIGH | 2h |
| 2.1.7 | DELETE /api/v1/services/:id | TODO | HIGH | 1h |
| 2.1.8 | Frontend: Service list page | TODO | HIGH | 3h |
| 2.1.9 | Frontend: Service form dialog | TODO | HIGH | 3h |

### 2.2 Customer Management
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 2.2.1 | Customer domain + repository | TODO | HIGH | 2h |
| 2.2.2 | Customer service layer | TODO | HIGH | 2h |
| 2.2.3 | GET /api/v1/customers | TODO | HIGH | 1h |
| 2.2.4 | GET /api/v1/customers/:id | TODO | HIGH | 1h |
| 2.2.5 | GET /api/v1/customers/search | TODO | HIGH | 2h |
| 2.2.6 | POST /api/v1/customers | TODO | HIGH | 2h |
| 2.2.7 | PUT /api/v1/customers/:id | TODO | HIGH | 2h |
| 2.2.8 | DELETE /api/v1/customers/:id | TODO | HIGH | 1h |
| 2.2.9 | Frontend: Customer list page | TODO | HIGH | 3h |
| 2.2.10 | Frontend: Customer form dialog | TODO | HIGH | 3h |
| 2.2.11 | Frontend: Customer search component | TODO | HIGH | 2h |

### 2.3 Orders Management
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 2.3.1 | Order domain + repository | TODO | HIGH | 3h |
| 2.3.2 | Order service layer | TODO | HIGH | 3h |
| 2.3.3 | Order items repository | TODO | HIGH | 2h |
| 2.3.4 | Order number generator | TODO | HIGH | 1h |
| 2.3.5 | GET /api/v1/orders | TODO | HIGH | 2h |
| 2.3.6 | GET /api/v1/orders/:id | TODO | HIGH | 2h |
| 2.3.7 | GET /api/v1/orders/:id/receipt | TODO | HIGH | 2h |
| 2.3.8 | POST /api/v1/orders | TODO | HIGH | 4h |
| 2.3.9 | PUT /api/v1/orders/:id | TODO | HIGH | 3h |
| 2.3.10 | PATCH /api/v1/orders/:id/status | TODO | HIGH | 3h |
| 2.3.11 | DELETE /api/v1/orders/:id | TODO | HIGH | 2h |
| 2.3.12 | Frontend: Order list page | TODO | HIGH | 4h |
| 2.3.13 | Frontend: Order form (new order) | TODO | HIGH | 6h |
| 2.3.14 | Frontend: Order detail page | TODO | HIGH | 3h |
| 2.3.15 | Frontend: Order status badges | TODO | HIGH | 2h |
| 2.3.16 | Frontend: Receipt print view | TODO | HIGH | 4h |

### 2.4 React Query Hooks
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 2.4.1 | use-orders hook | TODO | HIGH | 2h |
| 2.4.2 | use-customers hook | TODO | HIGH | 2h |
| 2.4.3 | use-services hook | TODO | HIGH | 2h |
| 2.4.4 | use-auth hook | TODO | HIGH | 2h |

### 2.5 Dashboard Layout
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 2.5.1 | Sidebar navigation | TODO | HIGH | 3h |
| 2.5.2 | Header component | TODO | HIGH | 2h |
| 2.5.3 | Dashboard layout wrapper | TODO | HIGH | 1h |
| 2.5.4 | Login page | TODO | HIGH | 3h |
| 2.5.5 | Register page | TODO | HIGH | 3h |

---

## Phase 3: Payments

### 3.1 Cash Payment
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 3.1.1 | Payment repository | TODO | HIGH | 2h |
| 3.1.2 | Payment service layer | TODO | HIGH | 2h |
| 3.1.3 | POST /api/v1/payments/cash | TODO | HIGH | 2h |
| 3.1.4 | Frontend: Payment dialog | TODO | HIGH | 3h |

### 3.2 Midtrans Integration
| # | Task | Status | Priority | MEDIUM |
|---|------|--------|----------|--------|
| 3.2.1 | Midtrans client package | TODO | MEDIUM | 2h |
| 3.2.2 | POST /api/v1/payments/midtrans/snap | TODO | MEDIUM | 3h |
| 3.2.3 | POST /api/v1/payments/midtrans/callback | TODO | MEDIUM | 4h |
| 3.2.4 | Frontend: QRIS payment flow | TODO | MEDIUM | 4h |

---

## Phase 4: Analytics

### 4.1 Analytics Backend
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 4.1.1 | Analytics service layer | TODO | MEDIUM | 3h |
| 4.1.2 | GET /api/v1/analytics/dashboard | TODO | MEDIUM | 3h |
| 4.1.3 | GET /api/v1/analytics/sales | TODO | MEDIUM | 2h |
| 4.1.4 | GET /api/v1/analytics/orders | TODO | MEDIUM | 2h |
| 4.1.5 | GET /api/v1/analytics/customers | TODO | MEDIUM | 2h |

### 4.2 Analytics Frontend
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 4.2.1 | Dashboard stats cards | TODO | MEDIUM | 3h |
| 4.2.2 | Sales chart (Tremor) | TODO | MEDIUM | 3h |
| 4.2.3 | Orders status chart | TODO | MEDIUM | 2h |
| 4.2.4 | Top customers table | TODO | MEDIUM | 2h |
| 4.2.5 | Analytics page layout | TODO | MEDIUM | 2h |

---

## Phase 5: Polish & Deploy

### 5.1 Polish
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 5.1.1 | UI/UX review + refinements | TODO | MEDIUM | 4h |
| 5.1.2 | Error handling + edge cases | TODO | HIGH | 4h |
| 5.1.3 | Loading states + skeleton | TODO | MEDIUM | 2h |
| 5.1.4 | Empty states | TODO | MEDIUM | 2h |
| 5.1.5 | Mobile responsiveness | TODO | MEDIUM | 4h |
| 5.1.6 | Zod validation (frontend forms) | TODO | HIGH | 4h |

### 5.2 Testing
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 5.2.1 | API endpoint testing (Postman/curl) | TODO | MEDIUM | 4h |
| 5.2.2 | User flow testing | TODO | MEDIUM | 4h |
| 5.2.3 | Edge case testing | TODO | MEDIUM | 2h |

### 5.3 Deployment
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 5.3.1 | Production environment setup | TODO | HIGH | 2h |
| 5.3.2 | Backend deployment (Railway/Render/VPS) | TODO | HIGH | 3h |
| 5.3.3 | Frontend deployment (Vercel) | TODO | HIGH | 2h |
| 5.3.4 | Domain + SSL setup | TODO | MEDIUM | 1h |
| 5.3.5 | CI/CD pipeline (GitHub Actions) | TODO | LOW | 4h |

### 5.4 Documentation
| # | Task | Status | Priority | Estimate |
|---|------|--------|----------|----------|
| 5.4.1 | README.md | TODO | LOW | 1h |
| 5.4.2 | API documentation | TODO | LOW | 2h |
| 5.4.3 | Update TASK.md | TODO | LOW | 30m |

---

## Time Estimate Summary

| Phase | Hours | Days (8h/day) |
|-------|-------|--------------|
| Phase 0: Setup | 20h | 2.5 days |
| Phase 1: Foundation | 45h | 5.5 days |
| Phase 2: Core Features | 90h | 11 days |
| Phase 3: Payments | 20h | 2.5 days |
| Phase 4: Analytics | 20h | 2.5 days |
| Phase 5: Polish & Deploy | 35h | 4.5 days |
| **TOTAL** | **230h** | **~29 days** |

---

## Current Sprint

### Week 1: Setup + Foundation
- [ ] 0.1.x - Infrastructure setup
- [ ] 0.2.x - Backend setup
- [ ] 0.3.x - Frontend setup
- [ ] 0.4.x - Migrations
- [ ] 1.1.x - API Response Package
- [ ] 1.2.x - JWT & Auth

---

## Notes

- Solo dev estimate: timeline flexible, < 1 month target
- Priority: HIGH = must have for MVP, MEDIUM = should have, LOW = nice to have
- Tasks will be marked DONE as implementation progresses
- Update this document weekly or when scope changes

---

*Document Version: 1.0*
*Last Updated: 2026-03-22*
