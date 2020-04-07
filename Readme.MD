# k8s-installed-features-catalogue

> What man is a man who does not make the world better.
>
> -- Balian, Kingdom of Heaven

## Abstract
A small definition for creating a catalogue of installed features on a
Kuberenetes cluster. Kubernetes is very flexible. Flexible to the extend that
you never know which features are enabled or installed.

## License
The license for the software is LGPL 3.0 or newer. 

## Usage
Just add the CRD to the cluster:

```bash
$ oc create -f crd-ionstalled-features-catalogue.yaml
customresourcedefinition.apiextensions.k8s.io/features.k8s.kaiserpfalz-edv.de created
$ 
```

And the you can create the information entries:
```bash
$ oc create -f ft-features-catalogue.yaml
installedfeature.k8s.kaiserpfalz-edv.de/features-catalogue created
$ 
```

Now you can select the info as cluster-admin. If you want to allow normal users to select it, you have to create a RBAC 
role for it:

```bash
$ oc create -f rbac-list-features.yaml 
clusterrole.authorization.openshift.io/listfeatures created
clusterrolebinding.authorization.openshift.io/feature:listfeatures created
$
```

## A note from the author
If someone is interested in getting it faster, we may team up. I'm open for that. But be warned: I want to do it 
_right_. So no short cuts to get faster. And be prepared for some basic discussions about the architecture or software 
design :-).

---
Bensheim, 2020-04-06-08:20+02:00
