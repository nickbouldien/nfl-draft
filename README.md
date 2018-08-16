# NFL DRAFT README

## Routes:
- / : index
- / : player/
- / : scouting
- / : test
- / : files/  (go to this route to see contents of the /files directory, or go to files/<file_name><file_extension>)


## start:
go run *.go

## Fetching data:
curl --request GET --url <base_url>/player/:id

ex: curl --request GET --url http://localhost:8080/player/3


## Drafting a player:
url: base url
id: id of the desired player
(make sure the header content-type is correct)

curl --request POST --url <base_url>/player/ --header 'content-type: application/x-www-form-urlencoded' --data id=3

ex: curl --request POST --url http://localhost:8080/player/ --header 'content-type: application/x-www-form-urlencoded' --data id=3


## postgres:
- sudo -u postgres psql
- CREATE DATABASE players_test1;
- \c players_test1
- CREATE TABLE ___;
- \d or \dt (describe - \d+ players)
- SELECT * FROM ____;
- UPDATE players SET drafted = true WHERE id = ___ RETURNING id;
- 
- players_test1=# SELECT position, count(position) FROM players GROUP BY position ORDER BY count DESC;


## Resources:
- ["Go in Practice"](https://www.manning.com/books/go-in-practice) by Matt Butcher and Matt Farina
- [Go docs](https://golang.org)
- [Simple JSON Rest API in Go (tutorial)](https://www.youtube.com/watch?v=hRR-Zy1H-Yo)
- Medium article ["Build RESTful API in Go and MongoDB"](https://github.com/mlabouardy/movies-restapi) by Mohamed Labouardy


## Cool packages I used/learned about:
- [goenv](https://github.com/joho/godotenv)


## Inspiration
This was inspired by [Key & Peele's](https://en.wikipedia.org/wiki/Key_%26_Peele) "East/West Bowl" comedy sketch(es), though this is the NFL draft (with one round) following the bowl
Players used from season 2 [video](https://www.youtube.com/watch?v=rT1nGjGM2p8)
Can find season 1 [video](http://www.cc.com/video-clips/5fndtz/key-and-peele-east-west-bowl) here

(disclaimer: the positions and year/class are made up by me (i.e., not explicity stated in the videos))

This is just a small thing so I could mess around with Go.  Feel free to correct/let me know of any of the (many) mistakes.


## TODO
- WRITE TESTS
- make an event table that records what team drafted which player (and when)
- allow user to update if player is drafted or not (undraft??)
- reset button to "undraft" all players and start over
- allow user to add more players to be drafted (from other seasons)
- write more idiomatic Go
