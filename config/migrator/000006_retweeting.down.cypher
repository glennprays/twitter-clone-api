MATCH (u:User { username: 'justbibir' })-[r:RETWEET]->(t:Tweet { content: 'This is my first tweet!' })
DELETE r