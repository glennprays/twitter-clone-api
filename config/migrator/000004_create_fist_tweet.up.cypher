MATCH (u:User { username: 'justbibir' })

CREATE (t:Tweet {
  content: 'This is my first tweet!',
  timestamp: datetime()
})

CREATE (u)-[:POSTED]->(t)