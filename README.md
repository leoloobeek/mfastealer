# mfastealer
This project includes sample code to help those wanting to use UiPath (RPA) to get around MFA/2FA VPN configurations. 

### What You Need
- UiPath Community Edition: https://www.uipath.com/community-edition-download    
- The mfastealer.exe binary from the Releases page, or if you don't trust me, Golang to compile the mfastealer.go file.    
- A webserver along with PHP
- A Windows SSH client for creating an SSH tunnel

### How's It Work?
Target is phished to enter credentials on a phishing page, including their MFA/2FA token(s). The phishing page captures the creds,
sends them down an SSH tunnel to an attacker's Windows system. The mfastealer.exe Go script works as an intermediary tool to
receive the captured crednetials and write them to a file. UiPath, which is a Robotic Process Automation (RPA) tool, will notice the
file was updated with credentials, and will open the VPN client, enter in all credentials, and initiate the VPN connection.

This entire process is fast enough (~5 seconds) to capture the credentials + token(s) and initiate the connection before the token(s) 
would typically expire. Obviously this is very dependent on the target infrastructure and VPN configuration.

### Setup
This is probably pretty new to most people, so if you have questions feel free to DM me on Twitter 
([@leoloobeek](https://twitter.com/leoloobeek)) or on the BloodHoundGang slack (@leo).

##### UiPath
Install the [UiPath Community Edition](https://www.uipath.com/community-edition-download). You can then import the `VPN Login.xaml` file to get started, this saved project
is for Anyconnect. I can't promise this will work for all Anyconnect clients and configurations, and definitely will not work
for other VPN clients. 

I'm absolutely no expert on UiPath and basically just fiddled with it until it worked for my needs. I would say it is pretty
intuitive and if you run into any issues hit me up. Few things to get you started:
- Import the provided `VPN Login.xaml` file to see what I used for the [video demo](https://vimeo.com/273241730)
- You will need to setup a `File change trigger` to watch when the `loginFile.txt` is written to. This file will be in the same
directory as mfastealer.exe. The mfastealer.exe tool will write username, password and token on each line of `loginFile.txt`.
- Either 'hardcode' or request the full path to the directory mfastealer.exe is in. This will allow you to run mfastealer.exe and
easily find the loginFile.txt that mfastealer.exe will write to
- Start up the VPN client either at the beginning with UiPath or yourself to save on time

As you'll see in the provided UiPath `VPN Login.xaml` file, I used a lot of `On Element Appear`, `Find Window` and `Attach Window` 
when interacting with the VPN GUI client. That seemed to work reliably for me, but there are plenty of ways to accomplish. You may
need to play around with it for 20 minutes or so to get the hang of it.

The provided demo UiPath project DOES NOT handle errors. A useful addition would be to detect if the VPN connection failed due to the
target fatfingering their password, then clearing out the loginFile.txt file awaiting for another target to fall for the attack.

##### mfastealer.exe
Next, download either the mfastealer.exe binary from the Releases page or compile it to run on Windows. Running the binary should
start a webserver listening on TCP port 3000. This webserver is expecting to receive an HTTP POST request with `username`, `password`,
and `secondPassword`. On the first set of credentials received, mfastealer.exe will write the credentials to loginFile.txt and
to cred.log. All subsequent credentials received will just be written to cred.log for safe keeping.

##### Phishing Page Setup
The attacker will setup a phishing page with a form requesting the username, password, and second password (token, "push", etc.).
The page should then send these credentials as `username`, `password`, and `secondPassword` respectively to the included
`login.php` page. Setup a local SSH port forward from your Windows system running mfastealer.exe and UiPath to the webserver. You
will need to forward the local system TCP port 3000 to the webserver's localhost:3000. Then the `login.php` file should post
credentials to localhost:3000, which should send the credentials down to mfastealer.exe. Once that happens, UiPath should pick up
the rest of the work!

### Special Thanks
Thanks to the guys over at FireEye for their ReelPhish [blog post](https://www.fireeye.com/blog/threat-research/2018/02/reelphish-real-time-two-factor-phishing-tool.html)
and [tool](https://github.com/fireeye/ReelPhish) release. The `login.php` file along with the overall idea of stealing tokens, came
directly from that tool.
