# Things API

CI Badges

Things API is an open-source, fast and scalable solution that enables you to speed up your IoT project development, by defining standard common features.

Things API rely on GoLang Gin Gonic web framework, MongoDB, Redis (for cache) and Sendgrid for mail management.

## Getting started
### Generate API Keys
To send mails (for user account management) things-api uses Sendgrid, so you should get an API key.

To resolve WiFi geolocations, things-api uses Google, so you should get an API key.

Create a `.env.prod` file from the included `.env` file, while customizing data such as domain name, API keys...

Install `docker` and `docker-compose`

Run `docker-compose up -d`

Watch `yourip:4000`, you should have a welcome message saying `Welcome to things API`.

Congratulations, you are all set !

## Sigfox Use
Connect to your account, select your device type and create two callbacks:
Type: `DATA` `UPLINK`,
URL: `https://youraddress/v1/sigfox/messages`
Method: `POST`,

Content-type: `application/json`

`mesType` is 1 for Sens'it, 2 for an Arduino Syntax, and 3 for Wisol 20

```
{
   	"sigfoxId":"{device}",
   	"frameNumber":{seqNumber},
   	"timestamp": {time},
   	"station": "{station}",
   	"snr": {snr},
   	"avgSnr": {avgSnr},
   	"rssi": {rssi},
    "mesType":1,
   	"data": "{data}"
}
```


Type: `SERVICE` `GEOLOC`,
URL: `https://youraddress/v1/sigfox/messages`
Method: `POST`,

Content-type: `application/json`

```
{
	"sigfoxId": "{device}",
	"timestamp": {time},
	"latitude": {lat},
	"longitude": {lng},
	"radius": {radius},
	"spotit": true
}
```

## NGinx configuration
### With HTTPS using certbot
You can copy paste and customize the [nginx/conf-https-step-1](https://github.com/IoThingsDev/api/tree/master/nginx/conf-https-step-1) to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload
Run certbot ....

Copy paste and customize the [nginx/conf-https-step-2](https://github.com/IoThingsDev/api/tree/master/nginx/conf-https-step-2)
to your `etc/nginx/sites-enabled/yourdomain`
`sudo service nginx reload`

###Without HTTPS

Copy paste and customize the conf-http
to your `etc/nginx/sites-enabled/yourdomain`
`sudo service nginx reload`


## API documentation and client generator
A swagger documentation is available [here](https://app.swaggerhub.com/apis/IoThings/Things-API/1.0.0), you can automatically generate a client in your favourite language!


## Roadmap
Some features would be nice to have, such as user roles management, Stripe billing management, Twilio SMS alerts.... And may be implemented in the future.

## Miscellaneous
If you want something you consider relevant to be implemented, feel free to fork the repo, and create a PR