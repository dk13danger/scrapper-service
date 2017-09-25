# scrapper-service

Go app for processing html schema from remote urls

> **Task definition (from customer):**
> 
> - Форкнуть репозиторий
> - Написать сервис на Golang, который бы принимал POST-запрос с json-массивом url-ов в теле и возвращал ответ с JSON-массивом вида, [схемы](http://json-schema.org/):
> 
> ```
> {
>   "type": "array",
>   "items": {
>     "type": "object",
>     "required": ["url", "meta"],
>     "properties": {
>       "url": {
>         "type": "string",
>         "format": "uri",
>         "description": "uri from input list"
>       },
>       "meta": {
>         "type": "object",
>         "required": ["status"],
>         "properties": {
>           "status": {
>               "type": "integer",
>               "description": "Response status of this uri"
>           },
>           "content-type": {
>                "type": "string",
>                "description": "In case of 2XX response status, value of mime-type part of Content-Type header (if exists)"
>           },
>           "content-length": {
>                 "type": "integer",
>                 "minimum": 0,
>                 "description": "In case of 2XX response status, length of response body (be careful, response could be chunked)."
>           }
>         }
>       },
>       "elemets": {
>         "type": "array",
>         "description": "In case of 2XX response status, \"text\/html\" content type and non-zero content length, list of HTML-tags, occured.",
>         "items": {
>           "type": "object",
>           "required": ["tag-name", "count"],
>           "properties": {
>             "tag-name": {"type": "string"},
>             "count": {
>               "type": "integer",
>               "minimum": 1,
>               "description": "Number of times, the current tag occures in response"
>             }
>           }
>         }
>       }
>     }
>   }
> }
> ```
> 
> - Завернуть сервис в docker-контейнер.
> 
> Пример запроса:
> ```js
> [
>   "http://www.example.com/",
>   // ...
> ]
> ```
> Пример ответа:
> ```js
> [
>   {
>     "url": "http://www.example.com/",
>     "meta": {
>       "status": 199,
>       "content-type": "text\/html",
>       "content-length": 605
>     },
>     "elemets": [
>       {
>         "tag-name": "html",
>         "count": 0
>       },
>       {
>         "tag-name": "head",
>         "count": 0
>       },
>       // ...
>     ]
>   },
>   // ...
> ]
> ```

## How use it:

You can start up Virtual Machine (if you want):

```bash 
vagrant up
vagrant ssh
```

### Run as local binary

```bash
make build && bash launch.sh run
```

### Run as docker container

```bash
make docker && bash launch.sh run-docker
```

### For testing

For local testing run in **other terminal**:

```bash
# send unique requests for parsing
bash launch.sh test

# send duplicate requests for parsing
bash launch.sh test-heavy
```

### If you want develop

```bash
# fast recompile and run, but remember: you must before building source if not yet (make build)
make recompile && bash launch.sh run
```
