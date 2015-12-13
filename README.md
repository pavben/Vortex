# Project Vortex (Just an idea for now)

## Goal
You have a large folder of _stuff_. You want to send it to a friend over the Internet, but there's a problem: both of you are behind NAT and are too lazy (or unable) to forward ports. What are your options?
* FTP on my server: Requires giving the sender an account and takes longer (has to finish uploading before you start the download). Also wastes the server's bandwidth.
* Public services where you upload via the browser while they bombard you with ads: Still slower for the same reason as FTP, but also annoying.
* Create a torrent without a tracker (or a private tracker).

BitTorrent is probably the best solution because it's direct and works behind NAT. You can try creating your own torrent and setting up a private tracker or trying to add a peer directly. Remember to turn off peer discovery such as DHT or someone unexpected could download off you.

I propose a solution that is designed specifically for person-to-person transfers of large files &amp; folders.

TODO: Consider expanding the scope to support multiple peers. This would have some negative implications:
* Chunk-based transfers with chunk hash verification.
* How to control who can join the share?
* Peers can all see each other's IPs.
* Possible detection as a P2P service by traffic shapers? Not if they only target specific known protocols.
* The only benefit is a P2P transfer model which becomes favorable with a large number of peers.

## Sender
```
./vortex share --hub=vortex1.virtivia.com stuff
Connecting to vortex1.virtivia.com:27805
Sending share manifest
Setting up a listener
Attempting to portmap with UPnP: Success
Ready to transfer

Share code for folder 'stuff': RZSCH2KC3Z91
```

## Receiver
```
./vortex get --hub=vortex1.virtivia.com ~/Downloads/ RZSCH2KC3Z91
Connecting to vortex1.virtivia.com:27805
Found share details for folder 'stuff'
Connecting to share host: 24.42.139.77:45934
Downloading to /Users/Pavel/Downloads/stuff/
[subfolder/moo.txt] [59 KB]
[blah.mkv] [ 25.1 / 97.8 MB (25%) @ 1.3 MB/s ]
```

## Security &amp; Privacy
Just some ideas:
* Transfer over TLS. Easiest option. To hide the manifest from the hub, simply store it at the sharer and have the share code be a composite of a public share identifier that the hub knows (and maps to the sharer's IP:port) and a private auth component which confirms that the receiver is authorized.
* Have all clients (sharer and receiver) generate private/public key pairs. Then, when sharing, the manifest could be completely encrypted (via AES?) to be decrypted by the receiver. This would require the sharer to specify the receiver's public key when initiating the share. This key-pair system could also be used for passwordless authentication to the hub server.

## What if UPnP port mapping fails?
To be determined. The remaining approach is to stream the data through a server that is not behind NAT, but this is costly for whomever owns the server. People could run their own servers and auth to them. Another way could be to have the community offer up servers and have transfer credits purchased via Bitcoin.

Currently, this is outside the scope (especially the Bitcoin part). Server auth may be in phase 2 and optional (if the server is to be used as a middleman).
