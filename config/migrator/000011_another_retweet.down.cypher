MATCH (u:User { username: 'justbibir' })-[r:RETWEET]->(t:Tweet { content: 'This is a reply tweet into reply!' })
DELETE r