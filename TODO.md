## TODO
	1. Get all playlist from spotify and display
		a. Be able to get the playlist id and use it to create playlist like melon top 100
## Commands

    SQLC Commands
    	- from the root of the directory, run: sqlc generate

    Goose Commands
    	- setup sqlc.yaml file
    	- inside the sql/schema directory, run: goose postgres postgres://postgres:@localhost:5432/melongo up

## PORTS 

	- broker-server
		a) http: 8080
	- auth-server
		a) http: 8081
		b) gRPC: 50001

## NOTES

	- Make a central database server where it's job is to soeley save it to the DB using message queue (non blocking)
		a) or maybe async grpc
	- Note on Auth Flow:
	  Just for now, using a manually generated api_key when user is created in the database. This is what
	  is being sent to the client side and client needs to include this with the subsequent requests.
	  When the request is made, it will direct to the auth-server and get the associated token for spotify.
	  Need to check if the accessToken is valid and if not, need to request a refreshToken. 
