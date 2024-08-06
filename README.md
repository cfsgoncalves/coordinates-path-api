
# Meight

A small API to manage trucks and orders and how to track orders that were shipped.

# Architecture
![alt text](image-1.png))

## Features

- Generates the best track for a truck based on orders assigned
- Generates events each time an order change a status
- Allows users to assigned orders to truck

## Not included

- Users permissions
- Authentification and Authentification ( expected to be shielded by an API gateway)

## Next steps

- Add prometheous metrics and Grafana for monitoring
- Add documentationt, like `swagger`

## Technologies used

- PostgresSQL for database
- Kafka for Message Queue
- Redis for Cache
- `SQLC` for generating sql queries
- `go-migrate` for migrating database schemas
- Golang as a programming language

## Run Locally

Install docker if you don't have it in your machine

Clone the project

```bash
  git clone https://github.com/cfsgoncalves/meight-project.git
```

Go to the project directory and run docker compose

```bash
  cd meight-project && docker compose up --build -d
```

## API Reference

#### Get all items

```http
  GET /api/items
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `api_key` | `string` | **Required**. Your API key |

#### Get item

```http
  GET /api/items/${id}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `string` | **Required**. Id of item to fetch |

#### add(num1, num2)

Takes two numbers and returns the sum.


