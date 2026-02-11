# Case Management System Backend

The Backend service for the Case Management System, built with **Go (Golang)** following **Clean Architecture** principles. This system features a dynamic form engine driven by metadata and a comprehensive workflow management for handling various types of business cases.

## 🚀 Key Features

- **Dynamic Form Engine**: Generates UI components based on database definitions (Metadata-driven).
- **Flexible Workflow**: Supports multi-step processes including Implementation, Approval, and Consulting.
- **Clean Architecture**: Decoupled layers for maintainability and scalability (Handler, UseCase, Domain, Infrastructure).
- **Audit Logging**: Comprehensive tracking of case transitions and attribute changes.
- **Multi-Environment Support**: Dedicated configurations for Dev, SIT, UAT, and Production.

---

## 🏗 Project Structure

This project follows the **Clean Architecture** pattern to ensure a separation of concerns:

```text
.
├── cmd/server/            # Entry points and Application bootstrap
├── configs/               # Environment-specific configuration files (.yaml)
├── internal/
│   ├── app/
│   │   ├── handler/       # HTTP Transport layer (Gin Gonic) and Routing
│   │   └── usecase/       # Business Logic and Orchestration
│   ├── domain/
│   │   ├── model/         # Domain entities and Data Transfer Objects (DTOs)
│   │   └── repository/    # Interfaces for data access
│   ├── infrastructure/
│   │   ├── database/      # GORM implementations for PostgreSQL
│   │   ├── api/           # External API Clients (Connectors)
│   │   └── seed/          # Initial master data seeding
│   └── pkg/               # Shared libraries (JWT, Responses, Utils)
├── migration/             # JSON files for data migration and definitions
└── utils/                 # General helper functions
```

## 🛠 Tech Stack

- **Language**: Go 1.23+
- **Framework**: Gin Gonic (HTTP Web Framework)
- **ORM**: GORM
- **Database**: PostgreSQL (with JSONB support)
- **Authentication**: JWT (JSON Web Token)
- **Containerization**: Docker & Docker Compose

---

## 📖 Core Concepts

### 1. Dynamic Attribute Engine

The system stores UI definitions in the database table `case_attribute_definitions`.


This allows the Frontend to render UI components dynamically  
(e.g. input, dropdown, grid) based on:

- `dataType`
- `config`

Values are provided by the API in a metadata-driven format.

---

### 2. ID + Label Data Contract

To reduce redundant API calls and improve performance, all selection-based data  
(e.g. dropdowns, multi-selects) follow a unified data structure:

```json
{
  "id": "UUID",
  "label": "Display Name"
}
```

This contract is consistently used across:
API Requests
API Responses
Database storage (JSONB)

### 3. State Machine & Workflow

Cases transition through predefined states based on:

- Action codes
- Role-based permissions
- Workflow configuration

State transitions are enforced by the business rules defined in the UseCase layer.

Example state transition:
```text
PENDING_ASSIGNMENT -> IN_PROGRESS -> WAITING_APPROVAL
```

Each transition may trigger:

- Audit log creation
- Case history tracking
- Permission re-evaluation
- Notification or downstream integration

---

## 🏃 Getting Started

### Prerequisites

- Docker
- Docker Compose
- Go 1.23+ (for local development)

---

### Running with Docker

1. Clone the repository.
2. Build and run the application:

```bash
docker-compose up --build
```

The backend service will be available at:
```bash
http://localhost:8080
```

---

### Local Development

Install dependencies:

```bash
go mod download
```
Run the application with environment selection:
```bash
# Example for Dev environment
go run cmd/server/main.go -env=dev
```

## 🔗 Internal Documentation

For detailed information on frontend integration and API contracts, refer to:

- **UI Component Catalog**
- **API Data Contract Specification**

These documents are crucial for correct Frontend–Backend integration.

---

## 🛡 Security

### Authentication

- JWT validation is required for most endpoints
- Token validation is handled via:
```bash
handler.ValidateToken()
```

---

### Authorization

- Role-Based Access Control (RBAC)
- Authorization rules are enforced at the **UseCase** layer


# UI Component Catalog & Data Contract

This document specifies the data structure and transfer protocols for the 13 standard UI components of ECM used in the Dynamic Form system.

## 1. Simple Data Types (Value as Primitive)
These components exchange values as simple strings, numbers, or booleans.

| DataType (BE) | UI Component | Value Format (JSON) | Note |
|---------------|--------------|---------------------|------|
| `string`      | Input (Text) | `"John Doe"`        | |
| `text`        | Textarea     | `"Long description"`| |
| `number`      | Input (Number)| `100`              | |
| `currency`    | Float        | `1234.50`           | |
| `radio`       | Radio Button | `"HOME"`            | Value is the selected `id` |
| `switch`      | Toggle       | `true`              | Boolean value |

## 2. Object Data Types (Value as Object)
Used when the UI needs to display a label but the database needs an ID.

| DataType (BE) | UI Component | Value Format (JSON) | Data Interaction |
|---------------|--------------|---------------------|------------------|
| `dropdownSingle` | Single Select | `{"id": "01", "label": "BKK"}` | **Transfer:** FE sends both ID and Label to BE to avoid re-fetching in Detail view. |
| `upload`      | File Upload  | `{"file_id": "uuid", "file_name": "a.jpg"}` | **Transfer:** After upload, FE sends the file metadata object. |

## 3. Array Data Types (Value as Array)

| DataType (BE) | UI Component | Value Format (JSON) | Data Interaction |
|---------------|--------------|---------------------|------------------|
| `dropdownMulti` | Multi Select | `[{"id": "1", "label": "A"}, {"id": "2", "label": "B"}]` | Array of ID/Label objects. |

---

## 4. Complex Component: List (Data Table / Grid)
The `list` component is a collection of other components arranged in a table.

### Data Schema
- **Configuration:** Defined in `config.columns`.
- **Value Storage:** A single JSON Array of Objects stored in a `jsonb` field.

### Column-Value Mapping Logic
Each object in the array uses the `code` from the column definition as the key.

| Column Type | Data Format within List Row | Example |
|-------------|----------------------------|---------|
| `string` / `number` | Primitive value | `"item-001"` |
| `dropdownSingle` | **Nested Object** | `{"id": "V-01", "label": "VISA"}` |

### Example Interaction (Dropdown inside List)
**Request/Response Payload:**
```json
{
  "def_attribute_id": "attr-grid-001",
  "dataType": "list",
  "value": [
    {
      "card_number": "4540-XXXX-XXXX-1111",
      "network": { "id": "net-01", "label": "VISA" }, 
      "reason": { "id": "res-01", "label": "Lost Card" }
    }
  ]
}