# POSTGRES USERLIST

This is a command line tool, that will list the users in a postgres database.

It will (optionally) also register a notifier/listener if the users table changes in any way.

## Rationale
Whenever a user is created or updated you may use this project to update the `/etc/pgpool2/pool_passwd` file.

## Installing

`go get cfeckardt/pg-userlist`

## Usage

`pg-userlist [-l] [-d delimeter] connection string`
