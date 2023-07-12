MATCH (u:User { username: 'justbibir' })-[p:POSTED]->(t1:Tweet { content: 'This is a quotation retweet!' })

MATCH (t1)-[qr:QUOTATION_RETWEET]->(t2:Tweet { content: 'This is my first tweet!' })

DELETE p, t1, qr