# User Auth Feature — Spec

### Requirement: User Registration

#### Scenario: Successful registration

- **WHEN** a user submits valid registration data
- **THEN** a new user account is created and a JWT token is returned
- **TEST** `TestUserRegistration_SuccessfulRegistration`

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

- **Decision:** JWT over session cookies — **Rationale:** Stateless design

## Risks / Trade-offs

- Token expiry management
