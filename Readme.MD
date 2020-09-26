# k8s-installed-features-catalogue

> What man is a man who does not make the world better.
>
> -- Balian, Kingdom of Heaven

## Abstract
A small definition for creating a catalogue of installed features on a
Kubernetes cluster. Kubernetes is very flexible. Flexible also means that it's hard to keep track of the activated or
installed features on a cluster.

Additionally, your cluster automation may check for the installed version of
features and act on that data (e.g. by updating a feature). Or you can check if
all required features for you are installed ...


## License
The license for the software is LGPL 3.0 or newer. 

## Usage

After the installation you can add the feature CRs to the k8s. 

Now you can get the features listed in a small nice manner:

```bash
$ oc get ift
NAME                 GROUP          VERSION        AGE       DOCUMENTATION
features-catalogue   cluster-info   1.0.0-alpha1   1d        https://github.com/klenkes74/k8s-installed-features-catalogue/
$
```

## A note from the author
If you want to get the end result faster, we may team up. I'm open for that. You have to keep in mind: I want to do it 
_right_. So no short cuts to get faster. Be prepared for some basic discussions about the architecture or software 
design :-).

---
Bensheim, 2020-04-06
