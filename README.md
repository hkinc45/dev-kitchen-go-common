# Dev Kitchen Go Common

This repository contains common Go libraries shared across the Dev Kitchen microservices ecosystem.

## Packages

### `auth`

The `auth` package provides middleware for Gin-based services to handle authentication and authorization.

#### Usage

1.  **Add to `go.mod`:**
    ```
    require github.com/hkinc45/dev-kitchen-go-common v0.2.0
    ```

2.  **Initialize the Middleware:**
    In your `main.go`, initialize the middleware. You will need to provide the OIDC provider URL and the service's own client ID.

    ```go
    import "github.com/hkinc45/dev-kitchen-go-common/auth"

    // ...

    authMiddleware, err := auth.NewMiddleware(
        context.Background(),
        os.Getenv("AUTH_PROVIDER_URL"),
        os.Getenv("OIDC_CLIENT_ID"),
    )
    if err != nil {
        log.Fatalf("Failed to create auth middleware: %v", err)
    }
    ```

3.  **Protect Routes:**
    You can now use the middleware to protect your Gin route groups.

    **For User-Facing Services:**
    Use `UserAuth()` to validate tokens from end-users. It checks the token signature, expiration, and ensures the service is in the token's audience (`aud` claim).

    ```go
    apiV1 := r.Group("/api/v1")
    apiV1.Use(authMiddleware.UserAuth())
    {
        // All routes in this group are now protected
        apiV1.GET("/me", ...)
    }
    ```

    **For Internal Services:**
    Use `ServiceAuth()` to validate service-to-service tokens. It checks the token signature, expiration, and ensures the token has the `internal-comm` role.

    ```go
    internalV1 := r.Group("/internal/v1")
    internalV1.Use(authMiddleware.ServiceAuth())
    {
        // All routes in this group are now protected
        internalV1.POST("/do-something", ...)
    }
    ```
