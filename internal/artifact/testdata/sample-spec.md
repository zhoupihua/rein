# User Auth Feature — Spec

## Context

The application needs user authentication to protect endpoints and personalize content.

## Goals

- Users can register and login
- JWT token-based session management
- Role-based access control

## Non-Goals

- OAuth/Social login (future phase)
- Passwordless auth

### Requirement: User Registration

#### Scenario: Successful registration

- **WHEN** a user submits valid registration data
- **THEN** a new user account is created and a JWT token is returned

#### Scenario: Duplicate email

- **WHEN** a user registers with an existing email
- **THEN** a 409 Conflict error is returned

### Requirement: User Login

#### Scenario: Successful login

- **WHEN** a user submits correct credentials
- **THEN** a JWT token is returned with 200 status

#### Scenario: Invalid credentials

- **WHEN** a user submits wrong password
- **THEN** a 401 Unauthorized error is returned

## Decisions

- JWT over session cookies for stateless design
- bcrypt for password hashing
