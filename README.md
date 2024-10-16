DNS snitch
==========

Sniff DNS responses and HTTP blocks size, and count bytes uploaded and downloaded, per domains.

Now you know why you should use [DoH](https://en.wikipedia.org/wiki/DNS_over_HTTPS).

Build it
--------

```
go build .
```

Run it
------

```
./dns-snitch eth0 eth1
```

Select sniffed interfaces. `sudo` can be required.
