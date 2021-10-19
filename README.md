# ibeacon simulator 

Simple ibeacon simulator written in golang.


Might need to run with `sudo` on rpi.

## Tests

For RPi portal beacon: 

```
EC9E84F8-87D8-498B-8B0C-9EF8D3AA94C7
```

Configure using [qrcode](https://npmjs.com/package/qrcode): 

```
qrcode EB4DAC16-7791-400A-A181-14F052F80B88
```

```
$ npm i -g qrcode
```


<!--
Notes:

Connect to BLE devices using RPi:
https://stackoverflow.com/questions/41707164/connect-ble-devices-with-raspberry-pi-3-b

How to allow non-root systemd service to use dbus for BLE operation
https://unix.stackexchange.com/questions/348441/how-to-allow-non-root-systemd-service-to-use-dbus-for-ble-operation/348460
>