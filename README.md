# prixa-assesment

# diagram
you can check the diagram [by click here](https://drive.google.com/file/d/1nvIuxDzK2PiLAQk9SZBq8aGv2Zj0Xd2Y/view?usp=sharing)

## How to Start Application

already make the docker-compose.yaml for the app so you can start the application easily.<br>
default port for api (:8000) <br>
default port for database (:6379) <br>

```bash
$ cd deployment
$ docker-compose up -d
```

# REST API

The REST API to the short url app is described below. The Default Port (8000)

## Post URL
`POST /short-url` => `localhost:8000/short-url` 
##### Body 
```
{
    "url": "www.google.com"
}
```

##### Response
```
{
    "data": {
        "url": "www.google.com",
        "short_url": "localhost:8000/ZGYwY2",
        "total_redirect": 0,
        "created_at": "2022-12-19T00:01:19.316281515Z",
        "updated_at": "2022-12-19T00:01:58.524020114Z"
    },
    "status": "success"
}
```

## Get short url
`GET /{url}` => `localhost:8000/j8dga5`
##### Response
it will be redirect to origin url


## Get Short Url status
`GET /{url}/status` => `localhost:8000/j8dga5/status`

##### Response
```
{
    "data": {
        "url": "www.google.com",
        "short_url": "localhost:8000/ZGYwY2",
        "total_redirect": 1,
        "created_at": "2022-12-19T00:01:19.316281515Z",
        "updated_at": "2022-12-18T00:01:58.524020114Z"
    },
    "status": "success"
}
```
