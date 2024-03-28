#!/bin/sh
cd src

go build -o goto .

./goto -http=:8080 -rpc=true &
master_pid=$!

./goto -master=127.0.0.1:8080 -http=:8081 &
slave_pid=$!

echo "Running master on :8080, slave on :8081."
echo "master pid: $master_pid, slave pid: $slave_pid"
echo "Visit: http://localhost:8081/add"
echo "Press enter to shut down"
read _

# Kill the master process if it is running
if [ ! -z "$master_pid" ]; then
    kill $master_pid
fi

# Kill the slave process if it is running
if [ ! -z "$slave_pid" ]; then
    kill $slave_pid
fi
