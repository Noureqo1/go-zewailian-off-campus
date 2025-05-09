# ADR 001: API Gateway Implementation

## 🧠 Context

In a microservices architecture, managing communication between frontend clients and multiple backend services can be complex. To simplify this, a **centralized entry point** is required. This approach allows us to handle:

- Request routing  
- Load balancing  
- Authentication  
- Cross-cutting concerns such as CORS and rate limiting  



## ✅ Decision

We implemented an **API Gateway** to serve as the single entry point for all client requests.

The gateway is responsible for:

-  Routing requests to the appropriate backend service  
-  Applying rate limiting and validating requests  
-  Handling CORS (Cross-Origin Resource Sharing)  
-  Performing basic authentication checks  

We selected **[KrakenD](https://www.krakend.io/)** due to its high performance and minimal configuration overhead.


## 🎯 Consequences

-  Centralized control over traffic and security  
-  Easier monitoring and logging of request flows  
-  Slightly increased complexity in gateway configuration  


