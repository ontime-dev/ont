#!/bin/bash

LOG_FILE="/var/log/ont.log"
CONF_DIR="/etc/ont"
CONF_FILE="${CONF_DIR}/ont.conf"
DBPASS_FILE="${CONF_DIR}/ont.pass"
DBPASS=""


#Create the mysql file from the following
create_ont_conf_file() {
	cat <<EOL > $CONF_FILE
SERVER_IP=''
SERVER_PORT='3033'
DEBUG='false'
EOL
	chmod 644 $CONF_FILE
}

store_dbpass_in_file() {
	echo "DBPASS=${DBPASS}" >> $DBPASS_FILE
	chmod 600 $DBPASS_FILE
}

get_ont_password() {
	read -s -p "Enter the new ontime DB password: " DBPASS
	echo  # Move to a new line after input

	read -s -p "Confirm password: " dbpass_confirm
	echo  # Move to a new line after input

	if [[ "$DBPASS" == "$dbpass_confirm" ]]; then
    		echo "Passwords match."
	else
    		echo "Passwords do not match."
	fi

	return $dbpass

}

create_db() {
	NEWPW=$1
	echo "Insert Mysql root password below"
	/bin/mysql -u root -p -e "CREATE DATABASE ontime; \n
	CREATE USER 'ont'@'localhost' IDENTIFIED BY '${NEWPW}'; \n
	GRANT ALL PRIVILEGES ON ontime.* TO 'ont'@'localhost';
	FLUSH PRIVILEGES;" 1>/dev/null
}

get_ont_password

create_db $DBPASS
echo "####Cleaned old setup and created new ontime DB####"

if [[ ! -d $CONF_DIR ]]; then 
	mkdir $CONF_DIR 
	chmod 755 $CONF_DIR
fi
echo "###${CONF_DIR} created###"

create_ont_conf_file
echo "Configuration file ${CONF_FILE} created"

store_dbpass_in_file
echo "###Stored new DBPass in ${DBPASS_FILE}"
