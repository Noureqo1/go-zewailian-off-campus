# Resilience Strategies

## Circuit Breakers
1. **Service-to-Service Communication**
   - Implementation: Using `gobreaker` package
   - Thresholds:
     - 5 failures in 10 seconds triggers open state
     - 30-second cool-down period
   - Applied to: All inter-service gRPC calls

2. **External API Calls**
   - Implementation: Custom circuit breaker with exponential backoff
   - Monitoring: Error rate and response time thresholds
   - Fallback: Cached responses where applicable

## Retry Mechanisms
1. **Database Operations**
   - Exponential backoff with jitter
   - Maximum 3 retries for transient failures
   - Deadlines to prevent infinite retries

2. **Message Publishing**
   - At-least-once delivery guarantee
   - Dead letter queue for failed messages
   - Retry with backoff for temporary failures

## Rate Limiting
1. **API Gateway**
   - Token bucket algorithm
   - Per-user and per-IP limits
   - Burst allowance for legitimate spikes

2. **WebSocket Connections**
   - Maximum messages per second per client
   - Connection pooling
   - Graceful degradation under load

## Fallback Strategies
1. **Chat Service**
   - Temporary message queue during Redis outage
   - Fallback to polling if WebSocket fails
   - Message persistence for offline delivery

2. **Resource Service**
   - Multiple storage regions
   - Cache frequently accessed resources
   - Progressive loading for large files

## Monitoring & Recovery
1. **Health Checks**
   - Regular service health probes
   - Database connection monitoring
   - WebSocket connection status

2. **Auto-Recovery**
   - Automatic service restarts
   - Data replication across zones
   - Automated failover procedures

## Justification
These strategies were chosen based on:
1. Real-time nature of the chat system
2. Need for high availability
3. Data consistency requirements
4. Resource efficiency
5. User experience preservation
