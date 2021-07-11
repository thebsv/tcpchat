## Simple TCP Server

### Requirements for the server

- Server room must run continually accepting new clients as they come in,
and removing them gracefully as they leave.

- Clients who connect to the server automatically join the server.

- Messages sent by each client should be broadcast to all other clients in the
server. This is the default action in the server which begins with a command '' 

- The clients must be able to list all other peers in the server using the command 'list'.
The server should send a response to the requesting client with the list of peers connected
to the current server.

- The clients can change their names in the server by using the command 'name <new_name>'. The server
should now display messages from this client with the new name. 

- The clients can exit the server gracefully by sending the command 'quit'.
