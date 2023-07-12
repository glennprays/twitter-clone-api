MATCH (u:User { username: 'justbibir' }), (t:Tweet { content: 'This is my first tweet!' })
CREATE (u)-[:RETWEET { timestamp: datetime() }]->(t)