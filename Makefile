define run
	@echo "Patching Runtime with ./$1/patch.diff..."
	@git apply ./$1/patch.diff
	@echo "Running $1 from source...\n"
	@cd $1 && ../go-src/bin/go run ./cmd
endef

default:
	@echo "Select make target..."

clean:
	@echo "Resetting Go to clean state..."
	@git restore go-src

build-go: clean
	@echo "Building Go from source..."
	@cd go-src/src && ./make.bash

patch:
	@echo "Creating patch..."
	@git diff go-src > patch.diff

run-hello-world: clean
	$(call run,1-test-patch)

jaeger:
	@docker stop jaeger > /dev/null || true
	@docker rm /jaeger > /dev/null || true
	@docker run -d --name jaeger \
      -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
      -e COLLECTOR_OTLP_ENABLED=true \
      -p 6831:6831/udp \
      -p 6832:6832/udp \
      -p 5778:5778 \
      -p 16686:16686 \
      -p 4317:4317 \
      -p 4318:4318 \
      -p 14250:14250 \
      -p 14268:14268 \
      -p 14269:14269 \
      -p 9411:9411 \
      jaegertracing/all-in-one:1.41
    $(shell open http://localhost:16686)

microservices: clean
	@cd example-app; \
	overmind start || true