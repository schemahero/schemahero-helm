#!/bin/sh

if [ $# -ne 1 ]; then
  echo "Usage: run.sh PGHOST"
  exit 1
fi

export PGHOST=$1
export PGPASSWORD='postgres'

cat manifest/db.yaml | sed "s/PGHOST/$PGHOST/g" | kubectl apply -f -

for i in $(seq 1 6); do
  kubectl get po -oname | grep 'testdb-controller' && break
  if [ $i -eq 6 ]; then
    echo "[$( date -Iseconds -u)] testdb controller hadn't started in 30s"
    kubectl logs 'sch-0'
    exit 1
  else
    echo "[$(date -Iseconds -u )] testdb controller not yet started, sleeping 5s"
    sleep 5
  fi
done

kubectl apply -f manifest/table.yaml
for i in $(seq 1 6); do
  psql -U postgres -t -c '\d' postgres | grep 'test_table' && break
  if [ $i -eq 6 ]; then
    echo "[$(date -Iseconds -u)] test_table hadn't been created in 30s"
    kubectl logs 'testdb-controller-0'
    exit 1
  else
    echo "[$(date -Iseconds -u)] test_table not yet created, sleeping 5s"
    sleep 5
  fi
done

