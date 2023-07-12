MATCH (u:User { username: 'justbibir' })-[p:POSTED]->(t1:Tweet { content: 'This is a reply tweet into reply!' })

MATCH (t1)-[r:REPLY]->(t2:Tweet { content: 'This is a reply tweet!' })

DELETE p, t1, r
