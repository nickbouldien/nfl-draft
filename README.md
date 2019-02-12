# NFL DRAFT README

## Inspiration
This was inspired by [Key & Peele's](https://en.wikipedia.org/wiki/Key_%26_Peele) "East/West Bowl" comedy sketch(es), though this is the NFL draft (with one round) following the bowl.
The players used are from season 2 [video](https://www.youtube.com/watch?v=rT1nGjGM2p8).
Season 1 video can be found [here](http://www.cc.com/video-clips/5fndtz/key-and-peele-east-west-bowl)

(disclaimer: the positions and year/class are made up by me (i.e., not explicity stated in the videos))

This is just a small thing I did as an excuse to mess around with Go for the first time

## Routes:
- `/index`
- `/players/`
- `/test`
- `/files/`  (use this route to see contents of the /files directory, or go to /files/<file_name><file_extension>)

## start:
`go run *.go`
or `go build` followed by `./nfl_draft`

## Fetch a player:
`curl --request GET --url <base_url>/players/:id`

ex: `curl --request GET --url http://localhost:8080/players/3`

## Fetch all players:
`curl --request GET --url <base_url>/players/`

ex: `curl --request GET --url http://localhost:8080/players/`

## Draft a player:
url: base url
id: id of the desired player

`curl --request POST --url <base_url>/players/id`

ex: `curl --request POST --url http://localhost:8080/players/3`

## TODO
- "appropriately" handle errors
- WRITE TESTS
- make an event table that records what team drafted which player (and when)
- allow user to add more players to be drafted (from other seasons)
- write more idiomatic Go
- swagger docs (or similar) for easy visibility of available routes
- DONE - reset button to "undraft" all players and start over

## postgres:
- `sudo -u postgres psql`
- `CREATE DATABASE players_dev;`
- `\c players_dev`
- `CREATE TABLE players;`
- insert the sample data (found in players.setup.sql)

## Resources:
- ["Go in Practice"](https://www.manning.com/books/go-in-practice) by Matt Butcher and Matt Farina
- [Go docs](https://golang.org)
- [Simple JSON Rest API in Go (tutorial)](https://www.youtube.com/watch?v=hRR-Zy1H-Yo)
- Medium article ["Build RESTful API in Go and MongoDB"](https://github.com/mlabouardy/movies-restapi) by Mohamed Labouardy
- [ENABLING CORS ON A GO WEB SERVER](https://flaviocopes.com/golang-enable-cors/), an article by Flavio Copes
- [How to not use an http-router in go](https://blog.merovius.de/2017/06/18/how-not-to-use-an-http-router.html), an article by Axel Wagner

## Cool packages I used/learned about:
- [goenv](https://github.com/joho/godotenv)
