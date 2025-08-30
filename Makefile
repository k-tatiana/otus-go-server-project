.PHONY: up-citus down-citus

# Запуск Citus кластера с 2 worker нодами
up-citus:
	docker-compose -f citus-docker-compose.yml -p citus up --scale worker=2 -d

# Остановка и удаление Citus кластера
down-citus:
	docker-compose -f citus-docker-compose.yml -p citus down

scale-citus:
	docker-compose -f citus-docker-compose.yml -p citus up --scale worker=4 -d

restart-citus:
	docker-compose -f citus-docker-compose.yml restart