# FJK Agent Architecture

## System Boundary

FJK Agent consists of three major runtime environments:

### Windows Desktop

- Tauri 2
- React
- TypeScript
- Desktop Pet
- Local control terminal

### Local Agent Core

- Go
- Eino
- Runtime
- Planner
- Executor
- Task
- Workflow
- Memory
- Permission
- Approval
- Event
- Conversation
- Tools
- SQLite

### Remote Services

- Relay Server
- Web Pet
- FJK Blog

## Core Rule

The desktop frontend must not directly execute privileged system operations.

All privileged operations must follow:

User
-> Task
-> Agent
-> Tool
-> Permission
-> Approval when required
-> Executor
-> Result
-> Event

## AI Boundary

Eino is responsible for:

- Agent reasoning
- Tool calling
- Graph execution
- Streaming
- Model integration

FJK Agent is responsible for:

- Task lifecycle
- Memory
- Permission
- Approval
- Workflow
- Audit
- Tool execution
- Desktop communication
- Blog management

## Security Boundary

Frontend must never directly access:

- AI API keys
- GitHub tokens
- local credentials
- unrestricted filesystem APIs
- unrestricted shell execution

Sensitive credentials will later be stored using Windows Credential Manager or equivalent secure storage.
