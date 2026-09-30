# Dev Kitchen Go Common: Package Reference

---

## 1. `auth`
- `NewMiddleware(ctx context.Context, providerURL, clientID, authServiceURL string) (*Middleware, error)`
- `(*Middleware) ServiceAuth() gin.HandlerFunc`: Protects internal routes using service-to-service credentials.
- `(*Middleware) RequireAuth() gin.HandlerFunc`: Enforces valid Bearer JWT.
- `RequirePermissionV2(httpClient *http.Client, resourceType string, idExtractor IDExtractor, action string) gin.HandlerFunc`: Fine-grained project and resource-level authorization.

---

## 2. `errors`
- `APIError`: Structured error type holding `StatusCode`, `Message`, and optional `Details`.
- `NewAPIError(code int, message string) *APIError`
- `NewNotFoundError(resource string) *APIError`
- `NewForbiddenError(message string) *APIError`
- `NewBadRequestError(message string) *APIError`

---

## 3. `resource_types`
Canonical constants for resource identifiers:
- `resource_types.Project = "project"`
- `resource_types.Recipe = "recipe"`
- `resource_types.StorageBucket = "storage-bucket"`
- `resource_types.VCSConnection = "vcs-connection"`
- `resource_types.Cluster = "cluster"`
- `resource_types.Region = "region"`

---

## 4. `models`
- `Recipe`: Core recipe definition, repository URLs, build arguments, attachments, and deployment status.
- `StorageBucket`: MinIO bucket properties, quotas, credentials.
- `RegistryToken`: Robot account credentials, expirations, permissions.
- `BillingAccount`: Billing balance, tier, payment methods.
- `ComputeProfile` & `RuntimePreset`: Hardware specifications and runtime configurations.
- `ObservabilityEvent`: Health checks, CPU/memory metrics, pod logs.

---

## 5. `metadata`
- `GenerateStandardLabels(projectID, recipeID, env string) map[string]string`: Returns standard Kubernetes labels (`app.kubernetes.io/name`, `app.kubernetes.io/instance`, `devkitchen.io/project`, etc.).
