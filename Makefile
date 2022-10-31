test-nats:
  #start nats server
	docker run --rm --name nats -p 4222:4222 nats:latest &
	sleep 1
	go test ./... -tags=nats -v -d
	docker stop -t 0 nats
