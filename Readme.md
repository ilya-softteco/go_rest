github.com/ilya-softteco/go_rest


docker run -v /home/dell/IdeaProjects/go_rest/restApiMySQLDB/migrations:/migrations --network host migrate/migrate -path=/migrations/ -database "mysql://user:password@tcp(localhost:3706)/db" up 1

