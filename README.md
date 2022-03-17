# apigo

## Technology used

- Golang - 1.17.8

## Curls

### /topsecret - POST

curl --location --request POST 'http://localhost:3000/topsecret' \
--header 'Content-Type: application/json' \
--data-raw '{
"satellites": [
{
"name": "kenobi",
"distance": 100.0,
"message": ["este", "", "", "mensaje", ""]
},
{
"name": "skywalker",
"distance": 115.5,
"message": ["", "es", "", "", "secreto"]
},
{
"name": "sato",
"distance": 142.7,
"message": ["este", "", "un", "", ""]
}
]
}'

### /topsecret_split/{satellite_name} - POST/GET

curl --location --request GET 'http://localhost:3000/topsecret_split/Kenobi' \
--header 'Content-Type: application/json' \
--data-raw '{
"distance": 200.0,
"message": ["este", "", "", "mensaje", ""]
}'
