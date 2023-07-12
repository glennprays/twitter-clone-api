MATCH (u:User { username: 'justbibir' })-[l:LIKES]->(t:Tweet { content: 'This is my first tweet!' })
DELETE l