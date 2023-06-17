#!/bin/sh

# Publish domain events
for i in {1..100}
do
	./gravity-cli pub accountCreated '{"id":'$i',"name":"fred'$i'","created_at":"2023-06-16T06:54:04.96Z"}'
done
