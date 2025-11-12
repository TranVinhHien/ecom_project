run:
	go run main.go
createredis:
	sudo docker run -d --name redis_c -p 6379:6379 -v /data/redis-data/:/data -e REDIS_ARGS="--requirepass 12345 --appendonly yes" redis:latest
dropredis:
	sudo docker rm redis_c
startredis:
	sudo docker start redis_c
stopredis:
	sudo docker stop redis_c
sqlc:
	sqlc generate
initmg:
	migrate create -ext sql -dir db/migration/ -seq init_mg
createc:
	sudo docker run --name mysql_c -e MYSQL_ROOT_PASSWORD=12345 -p 3306:3306 -d mysql:8.3.0
rmc:
	sudo docker rm mysql_c
initdb:
	sudo docker exec -it mysql_c mysql -u root -p12345 -e "CREATE DATABASE \`e-commerce\`;"
dropdb:
	sudo docker exec -it mysql_c mysql -u root -p12345 -e "DROP DATABASE \`e-commerce\`;"
startdb:
	sudo docker start mysql_c
stopdb:
	sudo docker stop mysql_c
buildimg:
	sudo docker build -t tranvinhhien1912/e-commerce:$(tag) 
pushimg:
	sudo docker push tranvinhhien1912/e-commerce:$(tag)
createtb:
	migrate -path db/migration/ -database "mysql://root:12345@tcp(localhost:3306)/e-commerce" -verbose up
droptb:
	migrate -path db/migration/ -database "mysql://root:12345@tcp(localhost:3306)/e-commerce" -verbose down
mountimg:
	sudo mount -t cifs //192.168.1.142/images_attr /home/master/project/e-commerce/imgs/products/images_attr   -o username=admin,password=Hienlazada#1,vers=3.0
	sudo mount -t cifs //192.168.1.142/images_customer  /home/master/project/e-commerce/imgs/avatar   -o username=admin,password=Hienlazada#1,vers=3.0
	sudo mount -t cifs //192.168.1.142/images  /home/master/project/e-commerce/imgs/products/images   -o username=admin,password=Hienlazada#1,vers=3.0 
#migrate -path db/migration/ -database "mysql://root:12345@tcp(localhost:3306)/e-commerce" goto 2 chuyển đổi phiên bảng
# ngrok http http://localhost:8080

.PHONY: run