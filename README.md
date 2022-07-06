# This is me - 02
Server with standard REST API methods (CRUD):
- Create.
- Get (Read through 'get' method).
- Update.
- Delete.
- List (Read through 'post' method).

### Third-party libraries
- "pq": Driver libraries for PostgreSQL to work with GO:
    - How to get:
        > `$ go get github.com/lib/pq`
- "ksuid": Libraries from 'Segmentio' to work with random generated strings
    - How to get:
        > `$ go get github.com/segmentio/ksuid`
- "jwt": Libraries from "Golang" (itshelf) to work with hashing (generating
JWT authorization tokens):
    - How to get:
        > `$ go get github.com/golang-jwt/jwt`
- "mux": Libraries from "Gorilla" to work with http handlers and endpoints:
    - How to get:
        > `$ go get github.com/gorilla/mux`
