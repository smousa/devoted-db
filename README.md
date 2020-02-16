# devoted-db
A simple in-memory database straight from stdin!

## Setup
This tool was written in [go1.13.6](https://golang.org/).  Build the binary by running `go build` and start the command-line prompt by typing `devoted-db`.

## Commands
| Operation  | Description                                                             |
|------------|-------------------------------------------------------------------------|
| `begin`    | Begins a new transaction.                                               |
| `commit`   | Commits all of the open transactions.                                   |
| `count`    | Returns the number of names that have the given value assigned to them. |
| `delete`   | Deletes the name from the database.                                     |
| `end`      | Exits the database.                                                     |
| `get`      | Prints the value from the given name.  Prints `NULL` if not present.    |
| `help`     | Help about any command.                                                 |
| `rollback` | Rolls back the most recent transaction.                                 |
| `set`      | Sets the name in the database to the given value.                       |
