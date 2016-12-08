all: clean build publish

clean:
	docker rmi agonzalezro/botella-test|true

build:
	docker build -t agonzalezro/botella-test .

publish:
	docker push agonzalezro/botella-test
