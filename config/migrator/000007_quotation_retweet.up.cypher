CREATE (t1:Tweet {
  content: 'This is a quotation retweet!',
  timestamp: datetime()
})
WITH t1
MATCH (u:User { username: 'justbibir' })
MATCH (t2:Tweet { content: 'This is my first tweet!' })

CREATE (u)-[:POSTED]->(t1)-[:QUOTATION_RETWEET]->(t2)