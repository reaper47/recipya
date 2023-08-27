---
sidebar_position: 3
---

# Nightly Build

A portable, standalone build is available on the [releases page](https://github.com/reaper47/recipya/releases/tag/nightly).
It is updated nightly.

First, download and extract the zip file compatible with your system. 
Please consult the [supported platforms](#supported-platforms) table if you are unsure which file to download.

Then, rename the `config.json.example` file to `config.json` and open it in a text editor 
to [modify the variables](/docs/installation/config-file). 

Finally, start the server by opening a command prompt in the folder, and run the following command.

```bash
./recipya-{os}-{architecture} serve
```

You can now access the website at the address specified.

:::note

Windows users need should replace http://0.0.0.0 with http://127.0.0.1 to access the website.

:::

The program cannot be updated within the interface yet. If you wish to update the app, please download the nightly
build and replace your existing `recipya` executable with the one from the build zip.
