# ADR 006: Redis Caching Layer

## 🧠 Context

To improve performance and reduce database load, we needed a fast and efficient caching mechanism for frequently accessed data such as:

-  Chat messages  
-  User sessions  
-  Auth tokens  


## ✅ Decision

We integrated **Redis** as a **caching layer** and **Pub/Sub** mechanism. It is used to:

-  Cache active sessions and recently accessed data  
-  Enable real-time messaging via Redis Pub/Sub channels  
-  Efficiently handle presence indicators in chat applications  


## 🎯 Consequences

-  Significantly improved response times and scalability  
-  Reduced load on the main PostgreSQL database  
-  Additional infrastructure and setup overhead  

