MATCH (u:User { username: 'justbibir' })-[p:POSTED]->(t1:Tweet { content: 'This is a reply tweet into another reply!' })

MATCH (u)-[l:LIKES]->(t1)
MATCH (t1)-[r:REPLY]->(t2:Tweet { content: 'This is a reply tweet into reply!' })
DELETE p, t1, r, l
