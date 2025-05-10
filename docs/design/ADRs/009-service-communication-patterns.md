# 9. Service Communication Patterns

Date: 2025-05-10

## Status

Accepted

## Context

Our chat application requires different services to communicate with each other:
1. Message Service - Handles message persistence and retrieval
2. WebSocket Service - Manages real-time connections and message broadcasting
3. User Service - Manages user authentication and profiles
4. Cache Service - Handles message and room caching

We need to decide on appropriate communication patterns between these services to ensure:
- Real-time message delivery
- System reliability
- Data consistency
- Scalability
- Error handling

## Decision

We have implemented a hybrid approach using both synchronous and asynchronous communication patterns:

### 1. Synchronous Communication (HTTP/REST)
Used for:
- User authentication and profile operations
- Room creation and management
- Message history retrieval
- Cache operations (with circuit breaker pattern)

Justification:
- Immediate consistency required for user operations
- Simple request-response pattern suitable for CRUD operations
- Direct feedback needed for user actions
- Circuit breaker pattern handles failures gracefully

### 2. Asynchronous Communication (WebSocket)
Used for:
- Real-time message broadcasting
- User presence updates
- Room activity notifications

Justification:
- Real-time updates required for chat functionality
- Reduces server load compared to polling
- Better user experience with instant message delivery
- Natural fit for event-driven architecture

### 3. Hybrid Patterns
Some operations use both patterns:
- Message sending:
  * Synchronous persistence to database
  * Asynchronous broadcasting to room participants
- Room updates:
  * Synchronous state change in database
  * Asynchronous notification to room members

## Implementation Details

1. **Synchronous Communication**:
```go
// HTTP handlers for synchronous operations
func (h *Handler) CreateRoom(c *gin.Context) {
    // Synchronous room creation with immediate response
}

func (h *Handler) GetMessages(c *gin.Context) {
    // Synchronous message retrieval with circuit breaker
}
```

2. **Asynchronous Communication**:
```go
// WebSocket handler for asynchronous operations
func (h *Handler) HandleWebSocket(c *gin.Context) {
    // Async message broadcasting
    hub.Broadcast <- message
    
    // Async user presence updates
    hub.Register <- client
}
```

3. **Resilience Patterns**:
```go
// Circuit breaker for synchronous operations
cbConfig := CircuitBreakerConfig{
    MaxRequests: 3,
    Interval:    10 * time.Second,
    Timeout:     30 * time.Second,
}

// Retry mechanism for transient failures
retryConfig := RetryConfig{
    MaxElapsedTime:   1 * time.Minute,
    MaxInterval:      5 * time.Second,
    InitialInterval: 100 * time.Millisecond,
}
```

## Consequences

### Positive
1. **Improved Reliability**:
   - Circuit breakers prevent cascade failures
   - Retry mechanisms handle transient issues
   - Asynchronous operations continue during partial outages

2. **Better Performance**:
   - Real-time updates via WebSocket reduce server load
   - Caching with circuit breakers improves response times
   - Async operations don't block the main request flow

3. **Enhanced User Experience**:
   - Immediate feedback for user actions
   - Real-time message delivery
   - Graceful degradation during issues

### Negative
1. **Increased Complexity**:
   - Managing both sync and async patterns
   - More complex error handling
   - Need for careful state management

2. **Eventual Consistency**:
   - Async operations may lead to temporary inconsistencies
   - Need for careful message ordering
   - Complex conflict resolution

## Monitoring and Metrics

To ensure the effectiveness of these patterns, we monitor:
1. WebSocket connection health
2. Circuit breaker state changes
3. Message delivery latency
4. Cache hit/miss rates
5. Error rates for sync/async operations

## Future Considerations

1. Consider implementing:
   - Message queuing for better reliability
   - Event sourcing for message history
   - Distributed tracing across sync/async operations

2. Potential improvements:
   - Implement server-sent events as fallback
   - Add message acknowledgment system
   - Enhance conflict resolution mechanisms
