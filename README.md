# FJK Agent

FJK Agent is a Windows local AI Agent platform with a desktop management terminal, embedded desktop pet, workflow system, tool system, memory, permissions, blog automation and remote Web Pet access.

## Project Status

Current stage:

- Phase 0 - Project Foundation

## Core Architecture

- Agent Core: Go
- Agent Framework: Eino
- Desktop: Tauri 2
- UI: React + TypeScript
- Database: SQLite
- Pet: Canvas
- Remote Relay: Go
- Blog: Astro + Markdown

## Main Modules

- Agent Core
- Runtime
- Planner
- Executor
- Task Engine
- Workflow Engine
- Memory
- Permission
- Approval
- Tool System
- Event System
- Conversation
- Desktop Terminal
- Pet System
- Relay
- FJK Blog integration

## Development Principles

1. Go Agent Core is the system brain.
2. Desktop UI is a control terminal, not the business logic layer.
3. Eino is the Agent execution framework, not the whole FJK Agent.
4. Tools must pass through the permission system.
5. High-risk actions require explicit approval.
6. Desktop and Web Pet share the same task system.
7. AI credentials must never be stored in frontend code.
8. Git Push and blog publishing must never bypass human approval.

## Development Phases

- Phase 0: Project Foundation
- Phase 1: Go Agent Core
- Phase 2: Desktop Terminal
- Phase 3: Pet Migration
- Phase 4: File / Git / GitHub Tools
- Phase 5: Long-term Memory
- Phase 6: Workflow Engine
- Phase 7: FJK Blog
- Phase 8: Relay + Web Pet
