# ADR 005: OAuth2 Authentication

## 🧠 Context

We needed a **secure**, **scalable**, and **user-friendly** authentication mechanism for both students and faculty members.  
A federated identity approach with **Single Sign-On (SSO)** was ideal to meet our usability and security goals.

## ✅ Decision

We implemented **OAuth 2.0** with **Google** as the authentication provider. This offers:

-  Simplified login experience  
-  No password storage or management  
-  Secure token-based authentication using **JWT (JSON Web Tokens)**  
-  JWTs manage user sessions and enable role-based access control for protected routes  

## 🎯 Consequences

-  Easier and safer user onboarding  
-  No internal password management required  
-  Full integration with JWT for access control  
-  Dependency on an external identity provider (Google)  

