---
title: Release Build
weight: 4
---

This section targets an installation without Docker.

The portable and standalone release builds are available on the [releases page](https://github.com/reaper47/recipya/releases) on GitHub.
The [nightly build](https://github.com/reaper47/recipya/releases/tag/nightly) is updated nightly if the main branch 
has new commits.

First, download the version of the software you wish to install compatible with your system, and extract the zip file. 
Please consult the [supported platforms](/guide/docs/installation/system-requirements) table if you are unsure which file to download.

Then, start the server by opening a command prompt in the folder, and run the following command.
The application will perform a one-time setup if not already done.

```bash
./recipya-{os}-{architecture} serve
```

You can now access the website at the address specified.

{{< callout type="info" >}}
The program cannot be updated within the interface yet. If you wish to update the app, please download the nightly
build and replace your existing `recipya` executable with the one from the build zip.
{{< /callout >}}

## Example

Let's say you have a Windows 11 computer, and you want to install Recipya v1.0.0 on it.

{{% steps %}}

### Access

You would first access the [releases page](https://github.com/reaper47/recipya/releases).

### Download

Identify version `v1.0.0` and download `recipya-windows-amd64.zip` under the **Assets** section.

### Extract

Extract the zip file on your computer.

### Run

Open a [command prompt](https://en.wikiversity.org/wiki/Command_Prompt/Open) and navigate to the folder you previously extracted.

```text
cd C:\path\to\recipya
```

Then, run Recipya once your command prompt is in its folder.

```text
.\recipya-{os}-{architecture} serve
```

### Enjoy

Open your browser to the address specified in the text of the command prompt.

If you see the following output:
```text
OK FDC database
OK Configuration file
Recipya is properly set up.
2023/12/27 19:08:06 goose: no migrations to run. current version: 20231204185640
Serving on http://127.0.0.1:8078
```

Then you would access `http://127.0.0.1:8078`.

{{% /steps %}}
