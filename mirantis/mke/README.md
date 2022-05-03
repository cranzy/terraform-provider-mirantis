# MKE Terraform Provider

This terraform provider integrates MKE into terraform via the MKE API.

The API provides nearly all functionality needed, although not all of it is
usabled.

## Implementations

### Provider

The provider is configured per MKE user, pointed at either a load balancer or
one of the MKE master nodes.

```
provider "mke" {
	endpoint = "https://${module.managers.lb_dns_name}"
	username = var.admin_username
	password = var.admin_password
}
```

### Data sources

#### ClientBundle

This Data Source retrieves a new client bundle from MKE.

The data source aims to provide sufficient information to allow configuration
of other providers such as the `kubernetes` provides.

```
data "mke_clientbundle" "admin" {}

provider "kubernetes" {
	host                   = data.mke_clientbundle.admin.kube[0].host
	client_certificate     = data.mke_clientbundle.admin.kube[0].client_cert
	client_key             = data.mke_clientbundle.admin.kube[0].client_key
	cluster_ca_certificate = data.mke_clientbundle.admin.kube[0].ca_certificate
}
```

The data source is still under development, and can be considered naive.
