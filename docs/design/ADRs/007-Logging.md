# ADR 004: Logging and Monitoring

## 🧠 Context

In a microservices environment with multiple services and Docker containers, **centralized logging** and **observability** are essential.  
These help with:

-  Debugging  
-  Performance monitoring  
-  Tracking system health  


## ✅ Decision

We adopted **OpenSearch** as our centralized logging solution.  
The setup includes:

-  Docker containers configured to ship logs directly to OpenSearch  
-  Logs are searchable and visualized via **OpenSearch Dashboards**  
-  Integrated alerting and metrics for **real-time monitoring**  


## 🎯 Consequences

-  Enhanced observability and deeper operational insights  
-  Faster detection and response to system issues  
-  Slight performance overhead due to continuous logging  
