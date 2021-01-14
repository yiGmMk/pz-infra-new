# -*- coding: utf-8 -*-

import MySQLdb, sys

#传入的参数
host = sys.argv[1]
user = sys.argv[2]
pwd = sys.argv[3]
db = sys.argv[4]

conn=MySQLdb.connect(host=host,user=user,passwd=pwd,port=3306)
cur=conn.cursor()

try:
    cur.execute('drop database if exists `' + db + '`')
except Exception as exc:
    print exc
cur.execute('create database `' + db + '` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci')

conn.commit()
cur.close()
conn.close()

