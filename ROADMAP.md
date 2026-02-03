# Roadmap: Go Common Library

## Phase 1: Epic: Centralized & Standardized Authorization Middleware

**Objective:** To create a single, powerful, and reusable authorization middleware that can be used by all microservices to enforce both project membership and fine-grained, IAM-style permissions.

*   **Milestone 1.1: Create Standard Middleware**
    *   `[ ] Task 1.1.1:` Create a new function `RequirePermission(permission string)` that returns a `gin.HandlerFunc`.
    *   `[ ] Task 1.1.2:` This middleware will extract the `projectId` from the URL parameters.
    *   `[ ] Task 1.1.3:` It will call the `auth-service`'s central `/internal/v1/auth/check` endpoint, passing the user's token, the project ID, and the required permission string.
    *   `[ ] Task 1.1.4:` If the `auth-service` returns "Permit", the request proceeds. If it returns "Deny", the middleware aborts the request with a `404 Not Found` to prevent leaking information.

*   **Milestone 1.2: Deprecate Old Middleware**
    *   `[ ] Task 1.2.1:` Mark any older, less specific permission-checking middleware functions as deprecated.
    *   `[ ] Task 1.2.2:` Update the library's documentation to guide developers to exclusively use the new `RequirePermission` middleware for all authorization needs.
