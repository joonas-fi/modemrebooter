[![Build Status](https://img.shields.io/travis/joonas-fi/modemrebooter.svg?style=for-the-badge)](https://travis-ci.org/joonas-fi/modemrebooter)
[![Download](https://img.shields.io/bintray/v/joonas/modemrebooter/main.svg?style=for-the-badge&label=Download)](https://bintray.com/joonas/modemrebooter/main/_latestVersion#files)

Sometimes modems just drop the internet connection in a way that the network comes back
only after being restarted. This is because modems usually are piles or garbage made out
of transistors.

This application periodically checks if the internet is down, and if down for enough time,
it reboots your modem to try to bring your funny cat videos back.

This could also be done universally (so no programming needed to support new types of
modems) with hardware (smart plug perhaps), and it'd probably be a lot more robust: what
if the internet is down and somehow the modem admin panel is also stuck so we cannot
reboot via software?


Usage
-----

Download binary for your OS/architecture (works for Raspberry Pi as well) combo from the
download link.

Configure it to automatically start at boot:

```
$ TODO
```


Supported garbage
-----------------

This application has "plugins" for different types of modems, currently:

| Model                                  | Code                     |
|----------------------------------------|--------------------------|
| TP-Link TL-MR6400 garbage              | tplinktlmr6400           |


How to build & develop
----------------------

[How to build & develop](https://github.com/function61/turbobob/blob/master/docs/external-how-to-build-and-dev.md)
(with Turbo Bob, our build tool). It's easy and simple!
