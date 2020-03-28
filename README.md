### Glee
[![Go Report Card](https://goreportcard.com/badge/github.com/tonyOreglia/glee)](https://goreportcard.com/report/github.com/tonyOreglia/glee)

## Play this engine at [tonycodes.com](https://tonycodes.com/) 

<iframe src='https://gfycat.com/ifr/ElatedBitterDragon' frameborder='0' scrolling='no' allowfullscreen width='640' height='396'></iframe>

<iframe src='https://gfycat.com/ifr/KeenDamagedAustraliancurlew' frameborder='0' scrolling='no' allowfullscreen width='640' height='404'></iframe><p><a href="https://gfycat.com/discover/cdc-gifs">from Cdc GIFs</a> <a href="https://gfycat.com/keendamagedaustraliancurlew-social-distancing-surgeon-general-coronavirus">via Gfycat</a></p>

### Overview 
Glee is a chess engine written in Golang. Glee adheres to the [Universal Chess Interface (UCI) protocol](https://en.wikipedia.org/wiki/Universal_Chess_Interface) over a websocket connection on port 8081. To connect with Glee from a websocket client, connect at `/uci`. 

Glee is a fully working chess engine, but is being continually improved. Progress can be tracked on the project page [here](https://github.com/users/tonyOreglia/projects/2) 


### Core Technical Concepts/Inspiration

Glee was originally undertaken as a method of learning the ins and outs of Golang. As a user, you may utilize specific packages for use in your own engine, use the code to learn about chess programming, use it as a backend to test your UCI frontend, or simply play the engine's command line interface. 

This engine is built using bitboard representation of the position. That is, a series of 64-bit unsigned integers are used to represent a given position and efficiently calculate legal moves via bitwise operations. 
An alpha-beta search algorithm is used to trim the potential moves tree.
A basic evaluation is used to evaluate a given position, based on a pawn value of +100 for white. 


### Getting Started/Requirements/Prerequisites/Dependencies
If you want to use the engine for any reason feel free to fork, download or clone the repo. 

To run the Websocket Server, run 
```
$ go run cmd/glee/main.go 
```

Note that the server will default to running in localhost on port 8081, if it should be run on a different IP Address you can override the value via the environment varialbe ADDR before starting the server. For example, 
```
$ export ADDR=157.230.180.254:8080
```

### Tests
Run 
```
go test ./...
```

### Contributing
Feel free to open a PR, I would be stoked. 

### TODO
See project [here](https://github.com/users/tonyOreglia/projects/2)

### Contact
- tony.oreglia@gmail.com
