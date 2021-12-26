# Web-cli builder
# web-cli is a command line, dispatcher api, web application where you
# can communicate with a server and execute your commands.
serve:
	clear
	echo "Building Application ..."
	go run cmd/main.go server

dispatch:
	clear
	echo "Importing Usages ..."
	go run cmd/main.go