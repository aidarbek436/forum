.SILENT:
dbrun: buildbackend
	docker run -dp 1999:1999 --name app back
buildbackend:
	docker build -t back .
dbstop: dbstop
	docker stop app
dbdelete:
	docker rm app
dbclear:
	docker rmi back