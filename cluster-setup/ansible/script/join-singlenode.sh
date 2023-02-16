## Master node's ip address
## please edit mater ip address
masterIP=

## Joining node's ip address
## please edit node ip address
nodeIP=

## k8s node name
## please edit node name
nodename=

## k8s node's role
## please edit node role
nodeRole=

## k8s node's type
## please edit node type
nodeType=

## Master's k8s token & hash
TOKEN=$(kubeadm token list | grep 'authentication' | cut -f 1 -d " ")
HASH=$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')

## Validation
## please edit hostname
if [ -z "$TOKEN" ]; then
	new_token=$(kubeadm token create)
	ssh root@$nodeIP "swapoff -a; kubeadm join $masterIP:6443 --token $new_token --discovery-token-ca-cert-hash sha256:$master_hash; exit"
else
	token=$TOKEN
	master_hash=$HASH
	ssh root@$nodeIP "swapoff -a; kubeadm join $masterIP:6443 --token $token --discovery-token-ca-cert-hash sha256:$master_hash; exit"
fi

## Checking node status
nodestatus=$(kubectl get nodes | grep $nodename | sed 's/   /,/g' | cut -d "," -f2)
while [ $nodestatus = "NotReady" ];
do
	nodestatus=$(kubectl get nodes | grep $nodename | sed 's/   /,/g' | cut -d "," -f2)
	nodestatus=$nodestatus
	echo -n .;
	sleep 1;
done

echo -e "\n"

## Labeling k8s node for roles
kubectl label nodes/$nodename kubernetes.io/role=$nodeRole --overwrite
kubectl label nodes/$nodename nodetype=$nodeType --overwrite

echo -e "\n"
