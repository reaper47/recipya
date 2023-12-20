---
weight: 4
---

# Nightly Build

A portable, standalone build is available on the [releases page](https://github.com/reaper47/recipya/releases/tag/nightly).
It is updated nightly if the main branch has new commits.

First, download and extract the zip file compatible with your system. 
Please consult the [supported platforms](#supported-platforms) table if you are unsure which file to download.

Then, start the server by opening a command prompt in the folder, and run the following command.
The application will perform a one-time setup if not already done.

```bash
./recipya-{os}-{architecture} serve
```

You can now access the website at the address specified.

:::note

The program cannot be updated within the interface yet. If you wish to update the app, please download the nightly
build and replace your existing `recipya` executable with the one from the build zip.

:::
