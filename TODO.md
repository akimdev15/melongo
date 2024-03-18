## TODO
- Note on Auth Flow: 
	Just for now, using a manually generated api_key when user is created in the database. This is what 
	is being sent to the client side and client needs to include this with the subsequent requests.
	When the request is made, it will direct to the auth-server and get the associated token for spotify. 
	Need to check if the accessToken is valid and if not, need to request a refreshToken. 
- TODO next: Need to map user to token and save it to the database. Also handle the case where when user access through /callabck, 
	only save the user if they don't exist. Otherwise, just get their api key and update the valid token

## Commands

    SQLC Commands
    	- from the root of the directory, run: sqlc generate

    Goose Commands
    	- setup sqlc.yaml file
    	- inside the sql/schema directory, run: goose postgres postgres://postgres:@localhost:5432/melongo up

