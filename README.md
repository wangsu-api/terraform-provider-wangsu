# Terraform Provider For wangsu

[![stars](https://img.shields.io/github/stars/wangsu-api/terraform-provider-wangsu)](https://img.shields.io/github/stars/wangsu-api/terraform-provider-wangsu)
[![Forks](https://img.shields.io/github/forks/wangsu-api/terraform-provider-wangsu)](https://img.shields.io/github/forks/wangsu-api/terraform-provider-wangsu)
[![Go Report Card](https://goreportcard.com/badge/github.com/wangsu-api/terraform-provider-wangsu)](https://goreportcard.com/report/github.com/wangsu-api/terraform-provider-wangsu)
[![Releases](https://img.shields.io/github/release/wangsu-api/terraform-provider-wangsu.svg?style=flat-square)](https://github.com/wangsu-api/terraform-provider-wangsu/releases)
[![License](https://img.shields.io/github/license/wangsu-api/terraform-provider-wangsu)](https://img.shields.io/github/license/wangsu-api/terraform-provider-wangsu)
[![Issues](https://img.shields.io/github/issues/wangsu-api/terraform-provider-wangsu)](https://img.shields.io/github/issues/wangsu-api/terraform-provider-wangsu)

<div>
  <p>
    <a href="https://www.wangsu.com">
        <img src="https://static-wcs.wangsu.com/portalnav/icon_1725418458199_836_logo-23_6_0_200.png" alt="logo" title="Terraform" height="200">
    </a>
    <br>
    <i>Wangsu Infrastructure for Terraform.</i>
    <br>
  </p>
</div>



* Tutorials: https://www.terraform.io

* [![Documentation](https://img.shields.io/badge/documentation-blue)](https://registry.terraform.io/providers/wangsu-api/wangsu/latest/docs)

* [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)

* Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

    

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.13.x
* [Go](https://golang.org/doc/install) 1.17.x (to build the provider plugin)

## Usage

### Build from source code

Clone repository to: `$GOPATH/src/github.com/wangsu-api/terraform-provider-wangsu`

```sh
$ mkdir -p $GOPATH/src/github.com/wangsu-api
$ cd $GOPATH/src/github.com/wangsu-api
$ git clone https://github.com/wangsu-api/terraform-provider-wangsu.git
$ cd terraform-provider-wangsu
$ go build .
```

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory,  run `terraform init` to initialize it.

### Configure proxy info (optional)

If you are beind a proxy, for example, in a corporate network, you must set the proxy environment variables correctly. For example:

```
export http_proxy=http://your-proxy-host:your-proxy-port  # This is just an example, use your real proxy settings!
export https_proxy=$http_proxy
export HTTP_PROXY=$http_proxy
export HTTPS_PROXY=$http_proxy
```

## Run demo

You can edit your own terraform configuration files. Learn examples from examples directory.

Now you can try your terraform demo:

```
terraform init
terraform plan
terraform apply
```

If you want to destroy the resource, make sure the instance is already in ``running`` status, otherwise the destroy might fail.

```
terraform destroy
```

## Developer Guide

### DEBUG

You will need to set an environment variable named ``TF_LOG``, for more info please refer to [Terraform official doc](https://www.terraform.io/docs/internals/debugging.html):

```
export debuggable=true
```

In your source file, import the standard package ``log`` and print the message such as:

```
log.Println("[DEBUG] the message and some import values: %v", importantValues)

```

### License

Terraform-Provider-Wangsu is under the Mozilla Public License 2.0. See the [LICENSE](LICENSE) file for details.