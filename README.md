# Dev Kitchen Go Common (`dev-kitchen-go-common`)

The shared Go library for all microservices in the **Dev Kitchen** cloud developer platform. `dev-kitchen-go-common` consolidates reusable authentication middleware, shared domain models, standardized error representations, HTTP client utilities, currency conversion, and Kubernetes metadata generators.

---

## Library Architecture & Modules

| Package | Purpose |
| :--- | :--- |
| `auth` | Gin middleware for OIDC JWT verification, service-to-service auth, token exchange, and V1/V2 permission checks (`RequirePermissionV2`). |
| `resource_types` | Canonical constants defining resource categories (`Project`, `Recipe`, `StorageBucket`, `VCSConnection`, `Cluster`, etc.). |
| `models` | Shared structs for recipes, storage buckets, container registry tokens, compute profiles, runtime presets, billing accounts, and observability events. |
| `errors` | Structured `APIError` formatting compliant with REST error semantics. |
| `clients` | Standardized HTTP response parsers and error wrappers for inter-service communication. |
| `gateways` | Gateway abstractions for routing and inter-service dispatch. |
| `worker` | Concurrent background worker pool implementation. |
| `currency` | Currency conversion routines supporting multi-region billing. |
| `metadata` | Helper functions for constructing standard Kubernetes labels and annotations. |

---

## Installation & Import

Import into any Go service:
```go
import (
    "github.com/hkinc45/dev-kitchen-go-common/auth"
    "github.com/hkinc45/dev-kitchen-go-common/errors"
    "github.com/hkinc45/dev-kitchen-go-common/models"
    "github.com/hkinc45/dev-kitchen-go-common/resource_types"
)
```

Ensure your `go.mod` references the module:
```go
module my-service

require github.com/hkinc45/dev-kitchen-go-common v0.0.0
```

---

## Quickstart: Securing a Gin Service

```go
package main

import (
    "context"
    "net/http"
    "os"

    "github.com/gin-gonic/gin"
    "github.com/hkinc45/dev-kitchen-go-common/auth"
)

func main() {
    r := gin.Default()

    authMiddleware, err := auth.NewMiddleware(
        context.Background(),
        os.Getenv("AUTH_PROVIDER_URL"),
        os.Getenv("OIDC_CLIENT_ID"),
        os.Getenv("AUTH_SERVICE_URL"),
    )
    if err != nil {
        panic(err)
    }

    // Public route
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // Internal service-to-service route
    internal := r.Group("/internal/v1")
    internal.Use(authMiddleware.ServiceAuth())
    {
        internal.GET("/data", func(c *gin.Context) {
            c.JSON(http.StatusOK, gin.H{"secret": "internal data"})
        })
    }

    r.Run(":8080")
}
```

---

## Documentation Navigation (Diátaxis)

- **Reference:**
  - [Package & Types Catalog](docs/packages.md)
- **How-to Guides:**
  - [Auth Middleware Integration Guide](docs/auth-middleware.md)
