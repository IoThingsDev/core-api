##Things API

CI Badges

Things API is an open-source, fast and scalable solution that enables you to speed up your IoT project development, by defining standard common features.

Things API rely on GoLang Gin Gonic web framework, MongoDB, Redis (for cache) and Sendgrid for mail management.

## Getting started
To send mails (for user account management) things-api uses sendgrid, so you should get an API key
Create a .env.prod file from the included .env file, while customizing data such as domain name, API keys...
Install docker and docker compose
Run docker-compose up -d

Watch yourip:4000, you should have a welcome message saying welcome to things API.
Congratulations, you are all set !

## NGinx configuration
### With HTTPS using certbot
You can copy paste and customize the nginx/conf-https-step-1 to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload
Run certbot ....

Copy paste and customize the conf-https-step-2
to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload

###Without HTTPS
Copy paste and customize the conf-http
to your etc/nginx/sites-enabled/yourdomain
sudo service nginx reload

## API Routes documentation
A swagger documentation is available [here]()

## Sigfox Use
Connect to your account, select your device type and create two callbacks:
Type: uplink data,
Method: post,
URL: https://youraddress/v1/sigfox/messages

Type: service geoloc,
Method: post,
URL: https://youraddress/v1/sigfox/locations

## Roadmap
Some features would be nice to have, such as user roles management, Stripe billing management, Twililo SMS alerts.... And may be implemented in the future.

## Miscellaneous
If you want something you consider relevant to be implemented, feel free to fork the repo, and create a PR