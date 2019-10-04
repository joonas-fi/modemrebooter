[![Build Status](https://img.shields.io/travis/joonas-fi/modemrebooter.svg?style=for-the-badge)](https://travis-ci.org/joonas-fi/modemrebooter)
[![Download](https://img.shields.io/bintray/v/joonas/dl/modemrebooter.svg?style=for-the-badge&label=Download)](https://bintray.com/joonas/dl/modemrebooter/_latestVersion#files)

Sometimes modems just drop the internet connection in a way that the network comes back
only after being restarted. This is because modems usually are piles or garbage made out
of transistors.

This application periodically checks if the internet is down, and if down for enough time,
it reboots your modem to try to bring your funny cat videos back.

This could also be done universally (so no programming needed to support new types of
modems) with hardware (smart plug perhaps), and it'd probably be a lot more robust: what
if the internet is down and somehow the modem admin panel is also stuck so we cannot
reboot via software?

Anyway, pure-software approach worked for me, but smartplug (or custom hardware) support
could be developed as a plugin in the future.


Usage
-----

Download binary for your OS/architecture (works for Raspberry Pi as well) combo from the
download link.

Write a `config.json` file (see [config.example.json](config.example.json)).

Configure it to automatically start at boot (you might need to run it with sudo):

```
$ ./modemrebooter write-systemd-unit-file
Wrote unit file to /etc/systemd/system/modemrebooter.service
Run to enable on boot & to start now:
        $ systemctl enable modemrebooter
        $ systemctl start modemrebooter
```


Supported garbage
-----------------

This application has "plugins" for different types of modems, currently:

| Model                                  | Plugin ID                |
|----------------------------------------|--------------------------|
| TP-Link TL-MR6400 garbage              | tplinktlmr6400           |
| ZyXEL VMG1312-B10D garbage             | zyxelvmg1312b10d         |


How to build & develop
----------------------

[How to build & develop](https://github.com/function61/turbobob/blob/master/docs/external-how-to-build-and-dev.md)
(with Turbo Bob, our build tool). It's easy and simple!
