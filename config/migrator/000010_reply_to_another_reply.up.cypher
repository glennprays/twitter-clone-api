MATCH (u:User { username: 'justbibir' })
CREATE (t1:Tweet {
  content: 'This is a reply tweet into another reply!',
  timestamp: datetime()
})
CREATE (u)-[:POSTED]->(t1)
WITH t1, u
MATCH (t2:Tweet { content: 'This is a reply tweet into reply!' })
CREATE (t1)-[:REPLY { timestamp: datetime() }]->(t2)
CREATE (u)-[:LIKES { timestamp: datetime() }]->(t1)