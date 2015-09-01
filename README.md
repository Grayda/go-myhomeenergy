go-myhomeenergy
===============

go-myhomeenergy is a library that connects to https://myhomeenergy.com.au/ and will eventually allow you to retrieve information about your electricity usage from your smart meter.

Usage
=====

Run `go get github.com/grayda/go-myhomeenergy`, import the library, then call

`myhomeenergy.Login("you@example.com", "password")`

If you're not comfortable with storing your password in plaintext, you can call LoginToken with a pre-calculated token:

`myhomeenergy.LoginToken("ABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCDEFABCDEF")`

To calculate your login token, use this pseudo-code:

`sha256(emailAddress + ":" + md5(password))`

Acknowledgements
================

Thanks to [Christina Porter][1] for her excellent PHP script. This code is being written from the ground up, but I referred to the script to work out the login hash

To-Do
=====

- [ ] Most of this library
- [ ] More robust error handling, including returning false on error

[1]: http://www.porters.co/2014/10/08/electricity-meter-data/
