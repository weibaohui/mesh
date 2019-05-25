rdns-server
========

The rdns-server implements the API interface of Dynamic DNS, its goal is to use a variety of DNS servers such as CoreDNS, Route53.

Currently we have implemented CoreDNS, CoreDNS will use etcd backend and rdns-server can do some CRUD operations on etcd.

The API doc can be found [here](https://github.com/rancher/rancher/wiki/Rancher-2.0-Dynamic-DNS-Controller#rancher-dynamic-dns-service)

## Building

`make`

## Running

For detailed steps, you can refer to [here](deploy/README.md)

## License
Copyright (c) 2014-2017 [Rancher Labs, Inc.](http://rancher.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
