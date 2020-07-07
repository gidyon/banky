init_mysql:
	@echo 'initializing mysql - Please wait ...'
	@docker run --rm -v=/Users/jessegitaka/Desktop/bank/deployments/compose/storage:/var/lib/mysql --name=mysql-init -e=MYSQL_ROOT_PASSWORD='hakty11' mysql:8.0.18

stop:	
	@docker rm -f mysql-init