build:
	docker build -t todo .
run:
	docker run -p 8084:8084 --rm -v $(pwd):/app -v /app/tmp --name todo-air todo
kubnginx:
	kubectl apply -f ./nginx/nginx-service.yaml -n default
nginx_clean:
	kubectl delete service nginx-service
	kubectl delete pod nginx-pod
