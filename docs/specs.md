# ipar-cli specs

## General features and anti-requirements

1. The following configuration items can be configured by the user:
1.1. The Archival Location (a directory on the user's system where where content
     is stored).
1.2. The amount of storage to be allocated for archival purposes.

1. Let user "donate" X number of MiB/GiB/TiB to the IPAR.

1.1. How the storage is being used is entirely up to the user.
1.2. No centralized mechanic that allows IPAR to use anyone's storage.
2. Let user control bandwidth to limit how much outgoing traffic can be
   allocated towards IPAR to not overwhelm user's Internet connection.

## Linux Integration

1. Allow IPAR to create a virtual network interface (eg. a "macvlan" device) so
   that IPAR can use traffic shapping tools (eg. "tc") to effectively enforce
   incoming and outgoing bandwidth limits.

   --- ipfs.netdev ---
   ```
   [Match]
   Name=ipfs0

   [Network]
   DHCP=yes
   ```

   --- ipfs.network ---
   ```
   [NetDev]
   Name=ipfs0
   Kind=macvlan

   [MACVLAN]
   Mode=bridge
   ```

   --- tc (10 mbit tx) ---
   ```
   tc qdisc add dev ipfs0 root tbf rate 10mbit burst 64kbit latency 50m
   ```

## Internet Archive Integration (archive.org)

### `ipar ia item `<url>`

1. The `ia item` command:
1.1. Verifies that the URL is valid (any 200 & non-404 response).
1.1.1. Example URL: https://archive.org/details/game-12537905211
1.2. Create a directory named after the item in the Archival Location.
1.2.1. Example directory name: `game-12537905211`
1.2. Downloads all the item's files to the new directory in the Archival Location.
1.2.1. The item directory is https://archive.org/download/<item> (can be derived
       from a hyperlink on the page: "SHOW ALL").
1.2.2. Only if there's both enough disk space and the size of the additional
       content doesn't exceed the allocated archival storage.
1.3. Pins the content using the `ipfs add`:
1.3.1. For files: `ipfs add --fscache --progress --pin --nocopy "$@"`
1.3.2. For directories: `ipfs add \
    --dereference-args \
    --fscache \
    --progress \
    --pin \
    --nocopy \
    --recursive \
    --wrap-with-directory "$@"` (for directories)
