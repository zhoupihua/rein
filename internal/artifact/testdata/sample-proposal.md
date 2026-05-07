# User Auth Feature — Proposal

## Why

The application needs user authentication to protect endpoints and personalize content. Currently there is no auth system.

## What Changes

Add user registration, login, and JWT-based session management.

## Goals

- Users can register and login
- JWT token-based session management
- Role-based access control

## Non-Goals

- OAuth/Social login (future phase)
- Passwordless auth

## Key Assumptions

- Users have email addresses
- JWT tokens are sufficient for session management
- bcrypt is acceptable for password hashing

## Open Questions

- Should we support refresh tokens?
- What is the token expiry policy?
