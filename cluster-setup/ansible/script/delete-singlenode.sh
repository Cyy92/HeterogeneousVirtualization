## Specify delete node
nodename=

kubectl drain $nodename --ignore-daemonsets --delete-local-data
kubectl delete node $nodename
