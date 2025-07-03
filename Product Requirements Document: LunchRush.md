# Product Requirements Document: LunchRush
## Collaborative Lunch Ordering Plugin for Huly

### Version 1.0

---

## 1. Overview

LunchRush is a collaborative Huly plugin that helps employees coordinate daily lunch orders efficiently by providing a centralized platform for restaurant selection, order management, and real-time coordination.

## 2. Problem Statement

Current lunch coordination suffers from:
- Scattered messages across multiple channels
- Missed orders and last-minute chaos
- No central visibility of who's ordering what
- Inefficient back-and-forth communication

## 3. Solution: "Trello meets lunch"

A collaborative experience inside Huly where team members can coordinate lunch orders in real-time with full visibility.

## 4. Core Requirements

### 4.1 Frontend - Huly Plugin

#### Lunch Session Display
- Show the current day's active lunch session
- Display session status (open/locked)
- Show countdown timer to lock time

#### User Participation
- **Join Session**: One-click to join the current lunch order
- **Meal Selection**: Select and specify meal choices
- **Live Participant View**: Real-time display of:
  - Who has joined
  - What each person is ordering
  - Current order totals

#### Restaurant Management
- **Propose**: Add restaurant suggestions
- **Vote**: Upvote/downvote restaurant options
- **Display**: Show voting results and selected restaurant

#### Order Coordination
- **Nominate Order Placer**: Select who will place the physical order
- **Lock Mechanism**: 
  - Set automatic lock time
  - Lock the session at specified time
  - Prevent changes after lock
- **Final Summary**: Display complete order details post-lock

#### Notifications
- Alert when session is about to lock
- Notify when order is locked
- Update when someone joins or changes their order

### 4.2 Backend - Go Microservice with Dapr

#### Required Dapr Building Blocks

**Pub/Sub Component**
- Real-time updates when users join
- Broadcast meal selections
- Notify all users of session state changes
- Push notifications for lock events

**State Store Component**
- Store current session data
- Maintain participant list and selections
- Track restaurant votes
- Persist order history
- Suggested: Redis for performance

**Optional Components**
- **Bindings**: Simulate restaurant API integrations
- **Secrets**: Store API keys for third-party services

#### Core Backend Functions
- Session creation and management
- Real-time state synchronization
- Vote tallying and restaurant selection
- Order compilation and locking logic
- WebSocket/SSE for live updates

## 5. Data Models

```go
type LunchSession struct {
    ID            string
    Date          time.Time
    Status        string // "open" | "locked"
    Restaurant    Restaurant
    Participants  []Participant
    OrderPlacer   *string
    LockTime      time.Time
}

type Restaurant struct {
    ID        string
    Name      string
    Votes     int
    ProposedBy string
}

type Participant struct {
    UserID    string
    Username  string
    MealChoice string
    JoinedAt  time.Time
}
```

## 6. User Workflows

### Creating and Joining
1. User creates or opens today's lunch session
2. System displays active session in plugin
3. Users click to join the session
4. Real-time updates show new participants

### Restaurant Selection
1. Users propose restaurant options
2. Participants vote on preferences
3. System displays live voting results
4. Top-voted restaurant is selected

### Order Management
1. Participants enter their meal choices
2. Live view shows everyone's selections
3. Someone volunteers as order placer
4. Session locks at specified time
5. Final summary displayed to order placer

## 7. Technical Requirements

### Performance
- Real-time updates with < 1 second latency
- Support 50+ concurrent users per session

### Architecture
- Clean, idiomatic Go code
- Proper Dapr component usage
- Clear separation of concerns
- RESTful API design

### UI/UX
- Clear, intuitive interface
- Responsive design
- Real-time visual feedback
- Accessible interaction patterns

## 8. Optional Bonus Features

If time permits, consider implementing:
- **Reorder**: "Order same as last week" functionality
- **Anonymous Voting**: Hide who voted for what
- **Scheduled Reminders**: Daily lunch notifications
- **Gamification**: 
  - Track frequent organizers
  - "Always late" badges
  - Order streak tracking

## 9. Deliverables

### Project Structure
```
/
├── plugin/          # Huly plugin frontend
├── microservice/    # Go + Dapr backend
└── README.md        # Setup instructions
```

### Documentation
- Clear setup instructions
- Architecture decisions explained
- Any assumptions or design choices

## 10. Evaluation Criteria

The solution will be evaluated on:
- **Collaborative Workflow Modeling**: How well the real-time features work
- **Code Quality**: Clean, idiomatic Go implementation
- **Dapr Usage**: Proper implementation of pub/sub and state management
- **User Interface**: Clear, intuitive plugin design
- **Overall Approach**: Thoughtfulness and clarity over polish

## 11. Timeline

- **Duration**: 1 week from challenge acceptance
- **Focus**: Core functionality over polish
- **Submission**: Fork and implement in designated folders