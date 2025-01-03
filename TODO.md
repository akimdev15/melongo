## TODO
	1. Need to prevent adding duplicate songs to the playlist 
	2. Try to write some tests to prevent a case where a song can't be found on spotify
	3. If album isn't discovered, maybe add a code to check to get the track. It sometimes works when searching for the track type
	4. Work on getting the top 100 list from melon and update playlist with the updated list.

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
