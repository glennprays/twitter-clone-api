MATCH (u:User { username: 'justbibir' })-[p:POSTED]->(t1:Tweet { content: 'This is a reply tweet!' })

MATCH (t1)-[r:REPLY]->(t2:Tweet { content: 'This is my first tweet!' })

DELETE p, t1, r
