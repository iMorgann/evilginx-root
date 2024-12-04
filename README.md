<p align="center">
    <img alt="Evilginx2 Title" src="https://raw.githubusercontent.com/kgretzky/evilginx2/master/media/img/evilginx2-title-black-512.png" height="60" />
    <img alt="Botguard" src="https://raw.githubusercontent.com/TomAbel/evilginx2/master/.github/bot.png?raw=true" height="60" />
    <img alt="yaomingface" src="https://raw.githubusercontent.com/TomAbel/evilginx2/master/.github/YaoMingFace.png?raw=true" height="100" />
</p>
<p align="center">
    <img alt="build_badge" src="https://img.shields.io/github/workflow/status/TomAbel/evilginx.botguard/Go?style=for-the-badge"/>
</p>

**evilginx2** is a man-in-the-middle attack framework used for phishing login credentials along with session cookies, which in turn allows to bypass 2-factor authentication protection.

This fork uses go-rod to retrieve a valid botguard token, which is then passed along to evilginx2 to get past the "Couldn't sign you in. This browser or app may not be secure" page

## Prerequisites

- Go 1.17
- OS capable of running Chrome

## Setup

```
wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb
sudo apt install ./google-chrome-stable_current_amd64.deb

git clone https://github.com/TomAbel/evilginx.botguard/
cd evilginx.botguard

go mod tidy
make

sudo ./bin/evilginx
```
## Possible issues

### "Sorry, something went wrong there. Try again" code 400
- use shorter domain name with valid TLD e.g. `phishing.com` (not `phishingdom.ain.sbk` nor `localhost.xyz`)
- install chromium browser on machine, guide [here](https://github.com/go-rod/go-rod.github.io/blob/master/compatibility.md)

### Submitting email takes too long

This is when go-rod is launched to retrieve the token. Some things can be done to speed up the process:

1. Remove or shorten SlowMotion time in code.
2. Preload browser so launching it takes less time.
3. Prefetch botguard tokens for specific emails. This only works when you know the victims.

## Video

[Demo](https://cdn.bunkr.to/6ZI7.mp4)

[Full setup](https://cdn.bunkr.to/VrCw.mp4)

## Contact

Want to say hello? Contact (base64) dG9ta3Jpc3RpYW5hYmVsQGdtYWlsLmNvbQ

Want to say thanks? https://ko-fi.com/tomabel

## License

**evilginx2** is made by Kuba Gretzky ([@mrgretzky](https://twitter.com/mrgretzky)) and it's released under GPL3 license. Code modifications not in the original evilginx2 repository but present here belong to me and are licensed under GPL3, unless stated otherwise.
