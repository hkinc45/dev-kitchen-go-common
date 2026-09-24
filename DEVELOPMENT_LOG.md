# Development Log

## Session: September 24, 2026

*   **Feature: Ephemeral Observability & Telemetry Shared Models (Milestone 28)**
    *   **Goal:** Establish canonical shared data models for live container log streaming, CI build logs, deployment warning events, instantaneous container metrics, and workload container discovery across services.
    *   **Implementation:**
        *   Created `models/observability.go`:
            *   `CIBuildLogStep`, `CIBuildLogResponse` for Gitea Actions step execution logs.
            *   `DeploymentIssueEvent` for ArgoCD sync failures and Kubernetes warning events.
            *   `ContainerMetricTelemetry` for real-time CPU millicores and memory MiB usage.
            *   `LogStreamRequest` for SSE streaming options (`Pod`, `Container`, `TailLines`, `Timestamps`, `Follow`, `Previous`).
            *   `RecipeWorkloadPod`, `RecipeContainerInfo` for multi-workload pod discovery and categorized container listing (`main`, `init`, `sidecar`) with restart counters.
        *   Updated `models/recipe.go` and `models/build.go` to include `LatestBuildJobID` and `LatestBuildRunID`.
    *   **Verification:** Verified Go tests with `go test ./...` passing cleanly.

## Session: September 22, 2026

*   **Feature: Universal Resource Attachment Data Model & Validation (Milestone 25)**
    *   **Goal:** Establish shared data structures, alias projection mappings, and contract validation for repository recipe resource attachments.
    *   **Implementation:**
        *   Extended `types/recipe.go` with `ResourceAttachment` struct, including fields `ID`, `RecipeID`, `AttachedRecipeID`, `AliasMapping`, and metadata timestamps.
        *   Defined validation methods for alias mapping keys, ensuring alphanumeric variable identifiers conforming to standard POSIX environment variable naming rules without collisions with reserved platform environment variables.
        *   Tagged and released `v0.5.27`.
    *   **Verification:** Ran `go test ./...` with 100% pass rate.

## Session: September 18, 2026

*   **Feature: Canonical Dev Kitchen Metadata & Label Standardization**
    *   **Goal:** Establish a single source of truth for platform workload metadata and eliminate legacy double-labeling.
    *   **Implementation:**
        *   Created package `metadata` (`metadata/labels.go`) exporting canonical constants:
            *   `LabelOwnerID` (`dev.kitchen/owner-id`)
            *   `LabelProjectName` (`dev.kitchen/project-name`)
            *   `LabelRecipeID` (`dev.kitchen/recipe-id`)
            *   `LabelRecipeName` (`dev.kitchen/recipe-name`)
            *   Matching `Annotation*` constants.
        *   Implemented `PlatformMetadata` struct with `StandardLabels()` and `StandardAnnotations()` helpers.
        *   Added `IsLegacyLabel` to detect and filter deprecated un-prefixed labels (`owner-id`, `project-name`, `recipe-id`, `recipe-name`).
        *   Added unit tests in `metadata/labels_test.go` with 100% pass rate.

## Session: October 23, 2025

*   **Feature: Improved Error Handling**
    *   **Goal:** To provide more specific error information from API clients.
    *   **Implementation:**
        *   Updated the generic `HandleResponse` function in the `/clients` package.
        *   The handler now checks for a `404 Not Found` HTTP status code and returns a specific, typed `NotFoundError` from the `/errors` package.
    *   **Impact:** Downstream services can now reliably check for 'not found' scenarios and handle them gracefully, making the entire system more robust and idempotent.

## Session: October 29, 2025

*   **Feature: Improved Error Handling**
    *   **Goal:** To provide more specific error information from API clients.
    *   **Implementation:**
        *   Updated the generic `HandleResponse` function in the `/clients` package.
        *   The handler now checks for a `404 Not Found` HTTP status code and returns a specific, typed `NotFoundError` from the `/errors` package.
    *   **Impact:** Downstream services can now reliably check for 'not found' scenarios and handle them gracefully, making the entire system more robust and idempotent.
## Session: May 3, 2026

*   **Feature: Technical Standardization & Verification**
    *   **Goal:** Complete the standardization baseline by adding unit tests and performance benchmarks for core library components.
    *   **Implementation:**
        *   Created `auth/middleware_v2_test.go` with mocked authentication service responses.
        *   Created `errors/errors_test.go` to verify centralized error handling middleware.
        *   Created `worker/worker_test.go` with concurrency benchmarks.
        *   Implemented a global log-suppression strategy for tests using `TestMain` and `io.Discard` to prevent terminal flooding during high-frequency benchmarks.
    *   **Impact:** Ensures architectural integrity and performance stability of the common library. Verified that the worker pool maintains a baseline of ~18µs per message processing overhead.

## Session: 2026-05-03 (Evening)

*   **Refactor: Logging Standardization**
    *   **Goal:** Align the common library with the team-wide standard of using `log/slog` for structured logging.
    *   **Implementation:**
        *   Migrated all legacy `log` calls to `log/slog` in `auth/middleware.go` and `clients/response_handler.go`.
        *   Updated version tag to `v0.5.12` to make the changes available to downstream services.
    *   **Impact:** Ensures consistent, structured logging across all services that use the common library, improving observability and debugging.
