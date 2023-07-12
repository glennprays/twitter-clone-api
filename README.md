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
### Docker
To start this project in docker:
1. Build the Docker Compose first
   ```
   docker compose build
   ```
2. Execute the Docker Compose useing 'up' command
   ```
   docker compose up
   ```
   this docker will run in port `80`

- Stopping Docker Compose
  ```
  docker compose down
  ```
- Migrate Database in Docker
  ```
  docker run -v {{ migration dir }}:/config/migrations --network host migrate/migrate -path=/config/migrations/ -database mysql://user:password@tcp(host:port)/dbname?query up
  ```
  ```
  migrate -database ${NEO4J_URL} -path db/migrations up
  ```