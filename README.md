### Server
The system has
1. "/login" handler for login
2. "/register" handler for registering user
3. "/posts" handler to manage posts

The post request goes through a handler to ensure that the username cookie is set. Otherwise it redirects to loginpage handler.

The modules folder contains files that implement both the user and post management.

They both implement ServerHTTP method to serve functions.