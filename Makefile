build:
	docker build -t todo .
run:
	docker run -p 8084:8084 --rm -v $(pwd):/app -v /app/tmp --name todo-air todo
