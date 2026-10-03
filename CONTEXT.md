# Project Blueprint: API Platform SaaS (B2D)

## 1. Project Overview
This project is a Business-to-Developer (B2D) REST API platform. It provides a developer dashboard to manage API keys, view audit logs, and monitor usage limits. The system separates the frontend (Dashboard) and backend (Core API). 

## 2. Tech Stack
- **Frontend:** Vue 3, Vite, TypeScript, Tailwind CSS. (Runs on port 5173).
- **Backend:** Golang, Fiber framework (v2). (Runs on port 8080).
- **Database:** MongoDB (via MongoDB Atlas) for all persistent data (Users, API Keys, Audit Logs).
- **Cache / Rate Limiting:** Redis (Upstash) for sliding-window rate limiting and Anti-DDoS at the application level.
- **Documentation:** Mintlify (Hosted externally, uses MDX, not integrated into the Golang codebase).
- **Workflow:** Monorepo structure. Both servers run concurrently using `make dev` (utilizing `air` for Go hot-reload and `vite` for Vue).

## 3. Core Architecture & Rules
- **API Key Management:** API Keys MUST be securely generated in Golang, hashed (e.g., using bcrypt or SHA-256) before saving to MongoDB. The plain text key is ONLY shown to the user once upon creation.
- **Audit Logging:** Every successful/failed request using an API key must be logged into MongoDB (Timestamp, Endpoint, Status Code, Latency) so it can be displayed in the Vue dashboard.
- **Rate Limiting:** Implemented via Fiber Middleware using Redis. Free users and Premium users have different request limits.
- **Authentication:** Standard JWT-based authentication for dashboard access, supplemented by OAuth (Google/GitHub).

## 4. Development Roadmap
- **Phase 1:** Monorepo Scaffolding (Vue setup, Golang Fiber setup, Makefile configuration).
- **Phase 2:** Database Connection & Schemas (MongoDB integration, defining User, API Key, and Audit Log models).
- **Phase 3:** Authentication System (JWT flow, manual login/register, Google/GitHub OAuth setup).
- **Phase 4:** API Key Lifecycle (Key generation, secure hashing, and REST endpoints for reveal/revoke/rotate).
- **Phase 5:** Security Middleware & Logging (Redis rate limiting, Free/Premium tier logic, active request audit logging).
- **Phase 6:** Frontend Dashboard UI (Vue dark mode interface, Pinia Auth state, data tables for Audit Logs and API Keys).
- **Phase 7:** Playground & Documentation (Custom Vue API playground for dashboard users, Mintlify MDX setup for public docs).

## AI Assistant Instructions:
- **Do not hallucinate files:** Always check this CONTEXT.md before proposing architectural changes.
- **Keep it modular:** Separate Golang routes, controllers, and middlewares cleanly. Use Vue composition API (`<script setup>`) for frontend.
- **Optimize tokens:** When executing a prompt, only edit the specific files requested by the user. Do not rewrite unmodified code.