#!/bin/bash
cd /root/go/src/github.com/dilfish/simple_cryptor
today=`date +"%Y-%m-%d"`
pass=""


function mysql_backup() {
	/usr/bin/mysqldump -u root -pLannister33 ${1} > ${1}.sql
	/bin/gzip ${1}.sql
	./simple_cryptor ${pass} ${1}.sql.gz 
	/bin/mv ${1}.sql.gz orig
	/bin/mv ${1}.sql.gz.se data
	/usr/bin/git add data/${1}.sql.gz.se
}


function filedir_backup() {
	/bin/cp -av ${2}/${1} .
	/bin/tar cvf ${1}.tar ${1}
	/bin/gzip ${1}.tar
	/bin/rm -rf ${1}
	./simple_cryptor ${pass} ${1}.tar.gz
	/bin/mv ${1}.tar.gz orig
	/bin/mv ${1}.tar.gz.se data
	/usr/bin/git add data/${1}.tar.gz.se
}


function file_backup() {
    /bin/cp ${2}/${1} .
    /bin/gzip ${1}
    ./simple_cryptor ${pass} ${1}.gz
    /bin/mv ${1}.gz orig
    /bin/mv ${1}.gz.se data
    /usr/bin/git add data/${1}.gz.se
}


mysql_backup mc
mysql_backup dilfish


filedir_backup letsencrypt /etc
filedir_backup conf /usr/local/nginx
filedir_backup etc /usr/local
file_backup fish.conf /root/go/src/github.com/dilfish/libsm/app/libsm


/usr/bin/git commit -m "backup-"${today}
/usr/bin/git push
