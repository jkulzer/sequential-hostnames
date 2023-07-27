# sequential-hostnames

## very much a todo project

This is a client-server pair.

The server is meant to be run on a machine supplying the hostnames to machines getting initialized.

The client is meant to run on first startup of the linux machine. It reaches out to the url specified with the `-url` flag where the server is listening.

The Server has a list of hostnames specified before starting the server. The file lives in the server folder and is called `hostnames`.

When the client reaches out to the server with the correct useragent (to prevent browsers doing weird things), the server returns the hostname given to the client in json and deletes that hostname out of the list.

Then the client gets the json and sets the hostname according to it.
