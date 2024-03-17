TODO 
1. when handling authorization, want to save user's info to the database and return both 
api_key which is generated (returned to the client) and the token received from spotify 
back to the broker service. Eventually, we want the broker service to store that in a redis cache
so that we don't have to look for the spotify api key everytime user makes a request. We will be sending both down to other microservices in case they need to make reqeust to the Spotify API. 


COMMANDS
1. goose postgres postgres://postgres:@localhost:5432/melongo_auth up
2. 
