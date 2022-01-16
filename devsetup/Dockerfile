FROM cassandra:3.11.6

ADD wait-for-it.sh /wait-for-it.sh

ADD init.cql /init.cql

RUN chmod +x /wait-for-it.sh

CMD cassandra -p pidfile -R; /wait-for-it.sh 127.0.0.1:9042 --timeout=0 -- cqlsh -f /init.cql && kill -9 $(cat pidfile) && sleep 2 && cassandra -f -R
