# simple_bank
Following this [tutorial]([url](https://www.youtube.com/watch?v=rx6CPDK_5mU&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&ab_channel=TECHSCHOOL))

docker images:
 -postgres: 14-lpine
 -pgadmin: pgadmin4 latest
 
## postgres:
	docker run --name practiceDB -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=admin -d postgres:14-alpine
## startdb:
	docker start practiceDB

## stopdb:
	docker stop practiceDB

## createdb:
	docker exec -it practiceDB createdb --username=root --owner=root simple_bank

## dropdb:
	docker exec -it practiceDB dropdb simple_bank
