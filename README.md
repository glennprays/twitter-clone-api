# Twitter Clone API

This project just a simple implementation of a social media platform similar to [Twitter](https://www.twitter.com).  
The API provides several functionalities for a single role, <b>User</b>, allowing them to perform the following actions:
- <b>Post a Tweet</b>: Users can create and share their own tweets with the community.
- <b>Like a Tweet</b>: Users have the ability to like or endorse a tweet posted by another user.
- <b>Re-Tweet a Tweet</b>: Users can share someone else's tweet on their own profile, giving credit to the original author.
- <b>Quotation Re-Tweet a Tweet</b>: Users can retweet a tweet while adding their own comment or perspective to it.
- <b>Reply to a Tweet</b>: Users can engage in conversations by replying to a specific tweet, creating threaded discussions.
   
The program uses [go-gin](https://github.com/gin-gonic/gin) as the framework, [neo4j](https://github.com/neo4j/neo4j) as the No-SQL Database (<i>Graph Database</i>), and using [go-migrate](https://github.com/golang-migrate/migrate) for database migrations.
  
It also employs various libraries:
- [neo4j-go-driver](https://github.com/neo4j/neo4j-go-driver)
- [GoDotEnv](https://github.com/joho/godotenv)
  
## Get Started
### Enviroment variables (.env)
Before starting the program, you need to set the `.env` file first:
1. Create `.env` file in the root directory
2. Copy the enviroment variables from `.env.example`
3. Fill the variables
### Docker
To start this project in docker:
1. Open Docker Directory at `/config/docker`
   ```
   cd config/docker
   ```
2. Build the Docker Compose first
   ```
   docker compose build
   ```
3. Execute the Docker Compose using 'up' command
   ```
   docker compose --env-file ./../../.env up
   ```
   this docker will run in port `80`

- Stopping Docker Compose
  ```
  docker compose down
  ```
### Database Migrations
Ensure that you have installed [go-migrate](https://github.com/golang-migrate/migrate). Before migrating the database, create a database in your MySQL.  
To run the database migrations:
- UP Migration
  ```
  migrate -database ${NEO4J_URL} -path db/migrations up
  ```
- DOWN Migration
  ```
  migrate -database ${NEO4J_URL} -path db/migrations down
  ```
> Note: in your local computer (without using docker) you need to add NEO4J_URL as enviroment variable
 ```
export NEO4J_URL="neo4j://user:password@host:port/"
 ```

 #### Migrate Database on Docker
1. Get golang app docker image id
   ```
   docker ps
   ```
2. Open the docker image command
   ```
   docker exec -it <image_id> /bin/bash
   ```
3. run migratons command