## TODO

- Make melon scrapper a separate repo/library so that I can use it my adding the module. Also create a REST API so that it can be used from different services as well.

- For microservice, create org in github to manage all repos

## Commands

    SQLC Commands
    	- from the root of the directory, run: sqlc generate

    Goose Commands
    	- setup sqlc.yaml file
    	- inside the sql/schema directory, run: goose postgres postgres://postgres:@localhost:5432/melongo up
