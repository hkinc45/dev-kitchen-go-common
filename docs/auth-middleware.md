# Dev Kitchen Go Common: Auth Middleware Integration Guide

This guide describes how to integrate `github.com/hkinc45/dev-kitchen-go-common/auth` into any Dev Kitchen microservice.

---

## 1. Initializing Middleware

```go
package main

import (
    "context"
    "os"
    "github.com/hkinc45/dev-kitchen-go-common/auth"
)

func initAuth() (*auth.Middleware, error) {
    return auth.NewMiddleware(
        context.Background(),
        os.Getenv("AUTH_PROVIDER_URL"), // e.g. http://localhost:8081
        os.Getenv("OIDC_CLIENT_ID"),    // e.g. dev-kitchen-api
        os.Getenv("AUTH_SERVICE_URL"),  // e.g. http://localhost:8081
    )
}
```

---

## 2. Protecting User Endpoints with Resource Permissions

Use `RequirePermissionV2` to check whether the requesting user has permission on a specific resource (e.g., recipe):

```go
recipeGroup := r.Group("/recipes")
recipeGroup.Use(authMiddleware.RequireAuth())

idExtractor := func(c *gin.Context) string {
    return c.Param("id")
}

recipeGroup.DELETE(
    "/:id",
    auth.RequirePermissionV2(
        http.DefaultClient,
        resource_types.Recipe,
        idExtractor,
        "recipe:delete",
    ),
    recipeHandler.DeleteRecipe,
)
```

---

## 3. Protecting Internal Microservice Routes

For routes that should only be accessible by trusted internal peers (Core, Recipe Service, Billing):

```go
internal := r.Group("/internal/v1")
internal.Use(authMiddleware.ServiceAuth())
{
    internal.GET("/recipes/:id", recipeHandler.InternalGetRecipe)
}
```
Inbound requests must pass either a valid shared service secret or an inter-service token.
