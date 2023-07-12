MATCH (u:User { username: 'justbibir' })
CREATE (t1:Tweet {
  content: 'This is a reply tweet into reply!',
  timestamp: datetime()
})
CREATE (u)-[:POSTED]->(t1)
WITH t1
MATCH (t2:Tweet { content: 'This is a reply tweet!' })
CREATE (t1)-[:REPLY { timestamp: datetime() }]->(t2)