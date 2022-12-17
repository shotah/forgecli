![Coverage](https://img.shields.io/badge/Coverage-65.2%25-yellow)
[![CI](https://github.com/shotah/forgecli/workflows/ContinuousIntegration/badge.svg)](https://github.com/shotah/forgecli/actions?query=workflow:ContinuousIntegration)
[![Update release version.](https://github.com/shotah/forgecli/workflows/PublishRelease/badge.svg)](https://github.com/shotah/forgecli/actions?query=workflow:PublishRelease)

[![Windows](https://img.shields.io/badge/Windows-0078D6?logo=windows&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Ubuntu](https://img.shields.io/badge/Ubuntu-E95420?logo=ubuntu&logoColor=white)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)
[![Mac OS](https://img.shields.io/badge/mac%20os-000000?logo=macos&logoColor=F0F0F0)](https://docs.github.com/en/actions/reference/workflow-syntax-for-github-actions#jobsjob_idruns-on)

# ForgeCLI

Package was created with the express intent to remove the guess work out of Mod acquiring and updating.

<!-- TOC -->

- [ForgeCLI](#forgecli)
  - [Parameters](#parameters)
  - [Installation](#installation)
    - [Windows](#windows)
      - [Chocolatey](#chocolatey)
      - [Scoop](#scoop)
    - [Mac](#mac)
      - [Homebrew](#homebrew)
  - [Forge API Key](#forge-api-key)
  - [Basic Usage](#basic-usage)
    - [Example Outputs](#example-outputs)
      - [Example Success](#example-success)
      - [Example Fail](#example-fail)
  - [Run with Client Install](#run-with-client-install)
  - [Json File](#json-file)
    - [Usage](#usage)
    - [Field Descriptions](#field-descriptions)
    - [Json Example](#json-example)
  - [Building and Testing](#building-and-testing)

<!-- /TOC -->

## Parameters

**NOTE:** Due to the lack of Version pinning this can lead to unexpected behavior if the publisher updates mod unexpectedly.

To get started make sure you have a [CursedForge API Key](https://docs.curseforge.com/#getting-started).

**Default Folder Locations**

- `Windows:` "%AppData%/Roaming/.minecraft/mods"
- `Mac:` "~/Library/Application Support/minecraft/mods"
- `Linux:` "~/Library/Application Support/minecraft/mods"

**Mod Basics**

- `Mod Release types:` Release, Beta, and Alpha.
- `Mod dependencies:` Required and Optional
- `Mod family:` Fabric, Forge, and Bukkit.
- `Mod MC Version:` 1.12.2, 1.18.2, etc.

**CLI Parameters**:

- `forgekey:` Required as a parameter, or as an ENV var of FORGEKEY and MODS_FORGEAPI_KEY
- `file:` Specially formatted json to manager larger sets of mods.
- `projects:` Project ids that can be easily obtained from the Forge itself.
- `destination:` Default is OS Client Mod folder, Target folder for the downloaded mods.
- `family:` Used to filter mods based on server type. Options are Forge, Fabric, and Bukkit
- `release:` Default is Release, Used to allow for Beta and Alpha mod downloads.
- `version:` Default is LATEST, but this is Minecraft VERSION. e.g. 1.18.2,
  - PARTIAL matching is enabled, e.g. use 1.18 to pull back 1.18.2, 1.18.1, 1.18 mods
- `clear:` Default is false, allows CLI to remove all mods before downloading new Mods.
- `client:` Default is false, if family is defined and client is defined, it will attempt to install the latest client.
- `dependencies:` Default is True, this uses the mods required dependencies to download missing mods.
- `debug:` Enable extra logging.

## Installation

### Windows

#### Chocolatey

Chocolatey based install
How do I install these tools?

[Install chocolatey url](https://docs.chocolatey.org/en-us/choco/setup)

Note: Currently waiting on **approval** by Chocolatey Admins

```bash
choco install forgecli --version=0.0.3 -Y
```

#### Scoop

Scoop Manifests
How do I install these tools?

[Install scoop url](https://scoop.sh/)

Add this bucket to scoop:

```bash
scoop bucket add shotah https://github.com/shotah/scoop-bucket.git
```

Install tools via scoop install:

```bash
scoop install forgecli
```

### Mac

#### Homebrew

Homebrew Manifests
How do I install these tools?

[Install Homebrew](https://docs.brew.sh/Installation)

Add this tap to brew:

```bash
brew tap shotah/tap
```

Install tools via brew install:

```bash
brew install forgecli
```

## Forge API Key

This is more complicated because you will be pulling/using the latest mod for the release of your game. To get started make sure you have a [CursedForge API Key](https://docs.curseforge.com/#getting-started). Then use it as a parameter for your build

## Basic Usage

Simple command to download latest fabric modules:

```bash
forgecli -forgekey '$2a$10...' -projects "416089,391366" -family "fabric" -debug
```

### Example Outputs

#### Example Success

```bash
time="2022-02-28T09:37:51-08:00" level=info msg="Starting Forge Mod lookup"
time="2022-02-28T09:37:51-08:00" level=info msg="Found Lastest FileID: 3667363 for Mod: 416089"
time="2022-02-28T09:37:51-08:00" level=info msg="Downloading: https://edge.forgecdn.net/files/3667/363/voicechat-fabric-1.18.2-2.2.24.jar"
time="2022-02-28T09:37:52-08:00" level=info msg="Files in Destination Folder:"
time="2022-02-28T09:37:52-08:00" level=info msg="  voicechat-fabric-1.18.2-2.2.24.jar  "
time="2022-02-28T09:37:52-08:00" level=info msg="Download Complete."
```

#### Example Fail

```bash
time="2022-02-28T09:52:14-08:00" level=info msg="Starting Forge Mod Lookup"
time="2022-02-28T09:52:14-08:00" level=error msg="could not find 391366 for minecraft version: 1.18.2 or family: fabric"
time="2022-02-28T09:52:14-08:00" level=error msg=Exiting...
```

## Run with Client Install

Includes `-client`, this tells the CLI to look at the `-family` parameter and find the version that corresponds to `-version` and download and install it for you.

Note: Currently only works with Fabric. I am not aware of Forge's version XML to query. If someone has more information I'd be happy for the help.

```bash
forgecli -forgekey '$2a$10...' -projects "416089,391366" -family "fabric" -client -debug
```

## Json File

### Usage

```bash
forgecli -forgekey '$2a$10...' -file "forgeMods.json" -family "fabric" -client -dependencies
```

### Field Descriptions

- `name:` is currently unused, but can be used to document each entry.
- `projectID:` is the id found on the CurseForge website for a particular mod
- `releaseType:` Type corresponds to forge's R, B, A icon for each file. Default Release, options are (release|beta|alpha).
- `fileName:` is used for version pinning if latest file will not work for you.
- `version:` used to override when Mojang releases minor updates and mods have not updated.
  - Does not currently apply to dependencies of Mod. e.g. Error if dependency does not support the main MC version.

### Json Example

```json
[
  {
    "name": "fabric api",
    "projectID": "306612",
    "releaseType": "release"
  },
  {
    "name": "fabric voice mod",
    "projectID": "416089",
    "releaseType": "beta",
    "version": "1.18.1"
  }
]
```

## Building and Testing

Make a `./.env` file in the root folder and add your forge key.

```text
FORGEKEY='$2a$10...'
```

Most common used commands can all be found in the [makefile](makefile).

Test command:

```bash
make test
```

Build command:

```bash
make build
# . ./forgecli should now be available to be used.
```
