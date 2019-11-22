# Golang Server

## Grab `set1-en_us.json`
If you're on Linux, run the script provided at `server/pull-data.sh`.

## Environment Variables

`PORT`: the port the web server listens to, defaults to 8080

## Routes

*After the very first game*:

Request:

```
POST /cards
[
	"01IO026",
	"01DE020"
]
```

Response:

```
{
  "id":"1197794836317474816",
  "href":"",
  "win":true,
  "result":{
    "Regions":["Ionia","Demacia"],
    "Set":{"Demacia":600,"Ionia":600}}
  }
```

*Updating Points after other games*:

Request:

```
PUT /cards/1197794836317474816
[
	"01IO026",
	"01DE020"
]
```

Response:

```
{
  "id":"1197794836317474816",
  "href":"",
  "win":true,
  "result":{
    "Regions":["Ionia","Demacia"],
    "Set":{"Demacia":1100,"Ionia":1100}}
  }
```

*Viewing Results*

Request:

```
GET /view/1197794836317474816
```

Response: a web page
