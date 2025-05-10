# ADR 008:  Microservice

## 🧠 Context

In the **Zewailian Off Campus** platform, real-time communication is a key feature. To support this, a **messaging microservice** is required with a relational data model to store:

- Active chat rooms  
- Chat messages between users  
- Metadata for querying and filtering  

The database schema needs to efficiently handle a high volume of messages and support features such as direct messages, group chats, and activity tracking.

## ✅ Decision

We designed and implemented a schema with two main tables: `rooms` and `messages`.

Key design choices include:

- **UUID-based primary keys** for scalability and distribution  
- **Foreign key constraints** to maintain referential integrity between messages and rooms  
- **Timestamps** to track message order and room activity  
- **Indexes** on frequently queried columns to optimize performance  

## 🎯 Consequences

- Supports **real-time messaging** with efficient query performance  
- Enables **filtering, sorting, and retrieval** of messages based on room or time  
- Easy to extend for future features (e.g., message reactions, file attachments)  
- Slight increase in storage due to indexing and metadata fields   

## SQL structure:

```sql
CREATE TABLE IF NOT EXISTS rooms (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner_id VARCHAR(36) NOT NULL,
    created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS messages (
    id VARCHAR(36) PRIMARY KEY,
    room_id VARCHAR(36) NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    username VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    type VARCHAR(50) NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    recipient VARCHAR(255),
    FOREIGN KEY (room_id) REFERENCES rooms(id)
);

CREATE INDEX IF NOT EXISTS idx_messages_room_id ON messages(room_id);
CREATE INDEX IF NOT EXISTS idx_messages_timestamp ON messages(timestamp);
CREATE INDEX IF NOT EXISTS idx_rooms_last_activity ON rooms(last_activity);
 



