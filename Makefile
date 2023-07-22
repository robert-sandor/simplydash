clean:
	rm -f ./build/simplydash

build: clean
	go build -o ./build/simplydash ./cmd/simplydash

run: clean build
	./build/simplydash -c ./local/config/config.yml -i ./local/images --log-level debug --access-logs

debug:
	dlv debug --headless ./cmd/simplydash -- -c ./local/config/config.yml -i ./local/images --log-level debug --access-logs

dockerBuild:
	docker buildx build . -t simplydash:dev

dockerRun: dockerBuild
	docker run --rm -p 8080:8080 \
		-v ./local/config:/app/config \
		-v ./local/images:/app/images \
		-v /var/run/docker.sock:/var/run/docker.sock \
		--name simplydash \
		simplydash:dev

