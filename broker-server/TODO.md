<h2>Overview of the Flow</h2>

UI --> Broker Service --> Auth Service (Synchronous Call) --> 
(Auth Service creates user, generates API key, adds to Redis, and returns API key to Broker Service) -->
Broker Service --> Player Service (Asynchronous Call) -->
(Player Service processes request and returns response to Broker Service) -->
Broker Service --> UI
