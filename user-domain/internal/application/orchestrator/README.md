# Orchestrator Package - Architecture & Dependencies

## 📋 Overview

Package `orchestrator` implements the **Saga Pattern** for distributed transaction management across microservices. It has been refactored following **SOLID principles** and **Clean Architecture** for better maintainability and extensibility.

## 🏗️ Architecture

### Dependency Structure

```
┌─────────────────────────────────────────┐
│         orchestrator.go                 │
│  (High-level orchestration logic)       │
│  ✓ Depends on INTERFACES only           │
└────────────┬────────────────────────────┘
             │
             ├──► workflow.Workflow (interface)
             ├──► outbound.Producer (interface)
             ├──► outbound.Subscriber (interface)
             ├──► outbound.Logger (interface)
             └──► outbound.WorkflowRunner (interface)
                  
┌─────────────────────────────────────────┐
│           workflow/                     │
│  - workflow.go (Workflow interface)     │
│  - updation_user.go (implementation)    │
│  ✓ Encapsulates saga state & logic      │
└────────────┬────────────────────────────┘
             │
             └──► command.* (interfaces)
             
┌─────────────────────────────────────────┐
│            action/                      │
│  - action_types.go (constants & types)  │
│  - user_update_actions.go               │
│  - payment_update_actions.go            │
│  ✓ Separated by domain/responsibility   │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│           command/                      │
│  - command.go (interfaces)              │
│  ✓ Defines contracts for saga steps     │
└─────────────────────────────────────────┘
```

## ✅ Key Improvements

### 1. **Dependency Inversion Principle (DIP)**
- **Before**: `orchestrator.go` depended on concrete type `*workflow.UserUpdationWorkflow`
- **After**: `orchestrator.go` depends on `workflow.Workflow` interface
- **Benefit**: Easy to add new workflow types without changing orchestrator

### 2. **Single Responsibility Principle (SRP)**
- **Before**: All actions in one 272-line file `updation_user.go`
- **After**: Separated into:
  - `action_types.go` - Constants and shared types
  - `user_update_actions.go` - User-related actions
  - `payment_update_actions.go` - Payment-related actions
- **Benefit**: Easier to maintain, test, and understand

### 3. **Open/Closed Principle (OCP)**
- **Before**: Hard to extend with new workflows
- **After**: New workflows just implement `Workflow` interface
- **Benefit**: Extensible without modifying existing code

### 4. **Interface Segregation**
- Clear separation of concerns through interfaces:
  - `command.Execution` - For executing commands
  - `command.Verification` - For verifying results
  - `command.Compensation` - For rollback operations
  - `command.Approval` - For final approval
- **Benefit**: Each component only depends on what it needs

## 📦 Package Structure

```
orchestrator/
├── orchestrator.go           # Main orchestrator (depends on interfaces)
├── noop_orchestrator.go      # No-op implementation for simple cases
├── README.md                 # This file
├── workflow/
│   ├── workflow.go           # Workflow interface definition
│   └── updation_user.go      # User update workflow implementation
├── action/
│   ├── action_types.go       # Shared types and constants
│   ├── user_update_actions.go    # User-related actions
│   └── payment_update_actions.go # Payment-related actions
└── command/
    └── command.go            # Command interfaces
```

## 🔄 Workflow Phases

The Saga workflow follows a **3-phase commit protocol**:

### Phase 1: Execute (Pending)
- Send pending requests to all services in parallel
- Each service prepares changes but doesn't commit
- If any fails → **Compensate** (rollback)

### Phase 2: Verify
- Wait for all services to acknowledge pending data
- Verify all services can accept the changes
- If any fails → **Compensate** (rollback)

### Phase 3: Approve (Commit)
- Send approval to all services in parallel
- All services commit their changes
- If any fails → **Compensate** (rollback)

### Compensation (Rollback)
- Executed if any phase fails
- Reverts all changes across all services
- Ensures data consistency

## 🚀 Usage

### Creating a New Workflow

1. **Implement the Workflow interface**:
```go
type MyWorkflow struct {
    state         SagaState
    executions    []command.Execution
    verifications []command.Verification
    compensations []command.Compensation
    approvals     []command.Approval
}

func (w *MyWorkflow) Run(ctx context.Context) error {
    // Implement saga logic
}

func (w *MyWorkflow) GetState() string {
    return w.state.String()
}

func (w *MyWorkflow) GetExecutionLogs() []*SagaExecuteLog {
    return w.executionLogs
}

func (w *MyWorkflow) GetLastError() error {
    return w.lastError
}
```

2. **Create actions in the action package**:
```go
// action/my_service_actions.go
type myServiceExecution struct {
    producer  outbound.Producer
    data      *MyData
}

func (c *myServiceExecution) Execute(ctx context.Context) error {
    // Send pending command
}

func NewMyServiceExecution(producer outbound.Producer, data *MyData) command.Execution {
    return &myServiceExecution{producer: producer, data: data}
}
```

3. **Wire it up in orchestrator**:
```go
func (o *orchestrator) ExecuteMyWorkflow(ctx context.Context, data *MyData) error {
    myWorkflow := workflow.NewMyWorkflow(o.producer, o.subscriber, data)
    return o.workflowRuner.Execute(ctx, myWorkflow.Run)
}
```

### Using the Orchestrator

```go
// In production with real infrastructure
orchestrator := orchestrator.NewWorkflowStarter(
    kafkaProducer,      // Producer implementation
    kafkaSubscriber,    // Subscriber implementation
    logger,             // Logger implementation
    temporalRunner,     // WorkflowRunner implementation (e.g., Temporal, Cadence)
)

err := orchestrator.ExecuteUserUpdation(ctx, oldUser, newUser)

// For simple cases without orchestration
orchestrator := orchestrator.NewNoopOrchestrator()
```

## 🧪 Testing

The new architecture makes testing much easier:

```go
// Mock the Workflow interface
type mockWorkflow struct {
    runFunc func(ctx context.Context) error
}

func (m *mockWorkflow) Run(ctx context.Context) error {
    return m.runFunc(ctx)
}

// Test orchestrator with mock
func TestOrchestrator(t *testing.T) {
    mockWf := &mockWorkflow{
        runFunc: func(ctx context.Context) error {
            return nil
        },
    }
    // Test with mock...
}
```

## 📊 Dependency Graph

```
Application Layer
    ↓ (depends on)
Domain Interfaces (outport/inport)
    ↓ (implements)
Infrastructure Layer
```

**Key Points**:
- ✅ High-level modules don't depend on low-level modules
- ✅ Both depend on abstractions (interfaces)
- ✅ Abstractions don't depend on details
- ✅ Details depend on abstractions

## 🔮 Future Extensions

Easy to add:
1. **New workflow types** - Just implement `Workflow` interface
2. **New services** - Add actions in `action/` package
3. **Monitoring** - Add decorators around workflow execution
4. **Retry logic** - Wrap workflow runner with retry mechanism
5. **Circuit breaker** - Add resilience patterns easily

## 📚 References

- **Saga Pattern**: [Microsoft Docs](https://docs.microsoft.com/en-us/azure/architecture/reference-architectures/saga/saga)
- **SOLID Principles**: [Uncle Bob Martin](https://blog.cleancoder.com/uncle-bob/2020/10/18/Solid-Relevance.html)
- **Clean Architecture**: [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

## 🎯 Benefits Summary

| Aspect | Before | After |
|--------|--------|-------|
| **Coupling** | Tight (concrete types) | Loose (interfaces) |
| **Testability** | Hard to mock | Easy to mock |
| **Extensibility** | Hard to add workflows | Easy - just implement interface |
| **Maintainability** | 272-line action file | Separated by concern |
| **Dependency Flow** | Inconsistent | Clean (DIP compliant) |

---

**Last Updated**: November 2025
**Maintainer**: Development Team

