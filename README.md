# NFL DRAFT README

## Routes:
- / : index
- / : player/
- / : test

## Fetching data:
curl --request GET --url <base_url>/player/:id

ex: curl --request GET --url http://localhost:8080/player/3


## Drafting a player:
url: base url
id: id of the desired player
(make sure the header content-type is correct)

curl --request POST --url <base_url>/player/ --header 'content-type: application/x-www-form-urlencoded' --data id=3

ex: curl --request POST --url http://localhost:8080/player/ --header 'content-type: application/x-www-form-urlencoded' --data id=3


## Recourses:
- ["Go in Practice"](https://www.manning.com/books/go-in-practice) by Matt Butcher and Matt Farina
- [Go docs](https://golang.org)
- [Simple JSON Rest API in Go (tutorial)](https://www.youtube.com/watch?v=hRR-Zy1H-Yo)
- Medium article ["Build RESTful API in Go and MongoDB"](https://github.com/mlabouardy/movies-restapi) by Mohamed Labouardy


## Cool packages I used/learned about:
- [goenv](https://github.com/joho/godotenv)


## Inspiration
This was inspired by Key & Peele's "East/West Bowl" comedy sketch(es), though this is the NFL draft (with one round) following the bowl
Players used from season 2 [video](https://www.youtube.com/watch?v=rT1nGjGM2p8)
Can find season 1 [video](http://www.cc.com/video-clips/5fndtz/key-and-peele-east-west-bowl) here

(the postions and year/class are made up by me (i.e., not explicity stated in the videos))