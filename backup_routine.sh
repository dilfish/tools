#!/bin/bash
cd /disk1/go/src/github.com/dilfish/simple_cryptor
today=`date +"%Y-%m-%d"`
pass=""


function mysql_backup() {
	/usr/bin/mysqldump -u root -pLannister33 ${1} > ${1}.sql
	/bin/gzip ${1}.sql
	./simple_cryptor ${pass} ${1}.sql.gz 
	/bin/mv ${1}.sql.gz orig
	/bin/mv ${1}.sql.gz.se data
	/usr/bin/git add data/${1}.sql.gz.se
	/usr/bin/git commit -m ${1}.${today}
}


function filedir_backup() {
	/bin/cp -av ${2}/${1} .
	/bin/tar cvf ${1}.tar ${1}
	/bin/gzip ${1}.tar
	./simple_cryptor ${pass} ${1}.tar.gz
	/bin/mv ${1}.tar.gz orig
	/bin/mv ${1}.tar.gz.se data
	/usr/bin/git add data/${1}.tar.gz.se
	/usr/bin/git commit -m ${1}.${today}
}


mysql_backup mc
mysql_backup dilfish


fildir_backup jj /disk1/cao
