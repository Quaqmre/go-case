## Requirements

* Go (1.18)
* Docker
* Docker Compose
* GNU make (optional)
##### **_Mongo Search_**

**URL** :   /fetch

**Method** : POST

**Request example**

```json
{
"startDate": "2016-01-21",
"endDate": "2016-03-02",
"minCount": 2900,
"maxCount": 3000
}
```

**Success Response**

**Code** : 200

**Response example**

```json
{
    "code": 0,
    "msg": "Success",
    "records": [
        {
            "key": "TAKwGc6Jr4i8Z487",
            "createdAt": "2017-01-28T01:22:14.398Z",
            "totalCount": 2800
        },
        {
            "key": "NAeQ8eX7e5TEg7oH",
            "createdAt": "2017-01-27T08:19:14.135Z",
            "totalCount": 2900
        },
        {
            "key": "NAeQ8eX7e5TEg9oH",
            "createdAt": "2017-01-27T08:19:14.135Z",
            "totalCount": 3100
        }
    ]
}
```
## Build & Run

With Docker Compose

```shell
cd <project dir>

# build
docker compose -f docker-compose.yml build

# run
docker compose -f docker-compose.yml up
# or run in background
docker compose -f docker-compose.yml up -d

# stop
docker compose -f docker-compose.yml down
```

With make

```shell
cd <project dir>

make build
make debug # run
make run # run in background
make stop
```