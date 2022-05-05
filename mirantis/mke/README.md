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

### Resources

#### ClientBundle

This resource retrieves a new client bundle from MKE.

The resource aims to provide sufficient information to allow configuration
of other providers such as the `kubernetes` provides.

```
resource "mke_clientbundle" "admin" {
	name = "admin" # this actually doesn't do anything, but TF needs at least one attribute.
}
```

This will give you enough data to configure some other providers such as kubernetes:

```
provider "kubernetes" {
	host                   = resource.mke_clientbundle.admin.kube[0].host
	client_certificate     = resource.mke_clientbundle.admin.kube[0].client_cert
	client_key             = resource.mke_clientbundle.admin.kube[0].client_key
	cluster_ca_certificate = resource.mke_clientbundle.admin.kube[0].ca_certificate
}
```

The resource is still under development, and can be considered naive:

# It doesn`t delete the client bundle, as it is not clear how to identify the key