MATCH (u:User { username: 'justbibir' })-[:POSTED]->(t:Tweet { content: 'This is my first tweet!' })
DETACH DELETE u, t