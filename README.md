# Critique

Open-source tool ( call it a microservice :wink: ) for managing customer feedbacks.
Critique **will be** able up and running in few minutes ( I hope ) , using the docker image.

API will have the following endpoints. Updating exising feedbacks and tags will also be available.

| Method | Endpoint | Description |
| ------ | :------- | :---------- |
| GET    | /v1/feedbacks     | get with pagination                      |
| GET    | /v1/feedbacks/:id | get feedback details                     |
| POST   | /v1/feedbacks     | post a feedback                          |
| GET    | /v1/tags          | get all tags(complaint, proposal, etc)   |
| POST   | /v1/tags          | create tags                              |