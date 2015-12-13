# Project Vortex (Just an idea for now)

## Goal
You have a large folder called _stuff_. You want to send it to a friend over the intertubes, but there's a problem: both of you are behind NAT and are too lazy (or unable) to forward ports. What are your options?
* FTP on my server: Requires giving the sender an account and takes longer (has to finish uploading before you start the download). Also wastes the server's bandwidth.
* Public services where you upload via the browser while they bombard you with ads: Still slower for the same reason as FTP, but also annoying.
* Create a torrent without a tracker (or a private tracker).

BitTorrent is probably the best solution because it's direct and works behind NAT. You can try creating your own torrent and setting up a private tracker or trying to add a peer directly. Remember to turn off peer discovery such as DHT or someone unexpected could download off you.

I propose a solution that is designed specifically for person-to-person transfers of large files &amp; folders, and just for fun, let's also design it with a focus on privacy where possible.

TODO: Consider expanding the scope to support multiple peers. It would have the following implications:
* Chunk-based transfers with chunk hash verification.
* How to control who can join the share?
* Peers can all see each other's IPs.
* Possible detection as a P2P service by traffic shapers? Not if they only target specific known protocols.
* The only benefit is a P2P transfer model which becomes favorable with a large number of peers.
* Too similar to BitTorrent.

## Sender
```
./vortex share stuff vortex1.virtivia.com
Setting up a listener
Attempting to portmap with UPnP: Success
Connecting to vortex1.virtivia.com:27805
Registering with the hub server
Got share code RZSCH2
Secret share key for folder 'stuff': RZSCH2-KC3Z91
Ready to transfer

Receiver command: ./vortex get [path] RZSCH2-KC3Z91 vortex1.virtivia.com

Current receivers:
==================
None
```

## Receiver
```
./vortex get ~/Downloads/ RZSCH2-KC3Z91 vortex1.virtivia.com
Connecting to vortex1.virtivia.com:27805
Found host for share code RZSCH2
Connecting to share host: 24.42.139.77:45934
Authorizing with key KC3Z91
Downloading to /Users/Pavel/Downloads/stuff/
[subfolder/moo.txt] [59 KB]
[blah.mkv] [25.1 / 97.8 MB (25%) @ 1.3 MB/s]
```

## Security &amp; Privacy
* Idea: Include the sharer's public key in the share key. This prevents man-in-the-middle attacks because the client will not be able to communicate with anyone who does not know the sharer's private key.
* Idea: Upon connecting to the sharer, the receiver provides its public key. The sharer then securely generates 256 bits which become the AES key for the remainder of the session, and sends this key encrypted via RSA using the receiver's public key.

## What if UPnP port mapping fails?
To be determined. The remaining approach is to stream the data through a server that is not behind NAT, but this is costly for whomever owns the server. People could run their own servers and auth to them. Server auth may be in phase 2 and optional (if the server is to be used as a middleman).
