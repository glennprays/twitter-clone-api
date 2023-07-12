MATCH (u:User { username: 'justbibir' }), (t:Tweet { content: 'This is a reply tweet into reply!' })
CREATE (u)-[:RETWEET { timestamp: datetime() }]->(t)