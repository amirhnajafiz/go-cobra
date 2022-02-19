# Web CLI

## What is this project?
Building a **Dispatching Server** using Golang programming language
and Golang Cobra package. 

You can create a command in a task and send it to the main sever, 
the server will execute that command for you by using a **Runner**. 

In this application you can create a task then you can run it
or change it or delete it. You can also see the list of the tasks
and their status.

## What is a Dispatching server?
The Web Dispatcher is the entry point for all external HTTP requests 
and the interface between all HTTP clients and the SAP system. 
The Web Dispatcher can work as load balancer for incoming requests 
which are distributed among all available application servers.

## Why using Cobra?
Cobra is a powerful library and tool in Golang that is used to
create CLI (Command-Line Interface) applications. 
Cobra does this by providing tools that automate
the process and provide key abstractions that increase
developer productivity.

## What are the project routes?
- `/tasks/{page?}` GET
  - page argument is optional 
  - returns a list of tags
- `/run` POST
  - running a task command
- `/tasks` POST 
  - creating a new task command
- `/task/{id}` DELETE
  - remove an existing task
- `/task/{id}` GET
  - returns a task status
- `/task/{id}` PUT
  - updates an existing task

## How to run the project?
After cloning the project:
```shell
git clone https://github.com/amirhnajafiz/Web-CLI.git
cd Web-CLI
```

Run the project:
```shell
make serve
```

If you need a guid, run:
```shell
make dispatch
```

## Project dependencies
- go 1.17
- cobra v1.3.0
- mux v1.8.0
- sqlite v1.2.6
- crypto v0
- ozzo-validation v3
