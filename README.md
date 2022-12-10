# ForgeCLI Guide

Package was created with the express intent to remove the guess work out of Mod acquiring and updating.

## Getting Your Forge API Key:

This is more complicated because you will be pulling/using the latest mod for the release of your game. To get started make sure you have a [CursedForge API Key](https://docs.curseforge.com/#getting-started). Then use it as a parameter for your build

## Basic Run:

Simple command to download latest fabric modules:

```bash
. ./forgecli.exe -forgekey '$2a$10...' -projects "416089,391366,552655" -family "fabric" -debug
```

## Example Outputs:

**Example Successful Run**

```bash
time="2022-02-28T09:37:51-08:00" level=info msg="Starting Forge Mod lookup"
time="2022-02-28T09:37:51-08:00" level=info msg="Found Lastest FileID: 3667363 for Mod: 416089"
time="2022-02-28T09:37:51-08:00" level=info msg="Downloading: https://edge.forgecdn.net/files/3667/363/voicechat-fabric-1.18.2-2.2.24.jar"
time="2022-02-28T09:37:52-08:00" level=info msg="Files in Destination Folder:"
time="2022-02-28T09:37:52-08:00" level=info msg="  voicechat-fabric-1.18.2-2.2.24.jar  "
time="2022-02-28T09:37:52-08:00" level=info msg="Download Complete."
```

**Example Failed Run**

```bash
time="2022-02-28T09:52:14-08:00" level=info msg="Starting Forge Mod Lookup"
time="2022-02-28T09:52:14-08:00" level=error msg="could not find 391366 for minecraft version: 1.18.2 or family: fabric"
time="2022-02-28T09:52:14-08:00" level=error msg=Exiting...
```

## ForgeCLI Usage:

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
- `dependencies:` Default is True, this uses the mods required dependencies to download missing mods.
- `debug:` Enable extra logging.

## JSON File usage:

**Basic Usage Command**

```bash
. ./forgecli.exe -forgekey '$2a$10...' -file "forgeMods.json" -family "fabric"
```

**Field Description**

- `name:` is currently unused, but can be used to document each entry.
- `projectID:` is the id found on the CurseForge website for a particular mod
- `releaseType:` Type corresponds to forge's R, B, A icon for each file. Default Release, options are (release|beta|alpha).
- `fileName:` is used for version pinning if latest file will not work for you.
- `version:` used to override when Mojang releases minor updates and mods have not updated.
  - Does not currently apply to dependencies of Mod. e.g. Error if dependency does not support the main MC version.

**Example json File Format**

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
  },
  {
    "name": "Biomes o plenty",
    "projectID": "220318",
    "fileName": "BiomesOPlenty-1.18.2-15.0.0.100-universal.jar",
    "releaseType": "release"
  }
]
```

## Manually Building and Testing:

Make a `./.env` file in the root folder and add your forge key.

```text
FORGEKEY='$2a$10...'
```

CMD To test

```bash
go test ./...
```

CMD To Build

```bash
go build
# . ./forgecli should now be available to be used.
```

## TODO List

- Add fabric xml download so we can parse latest version
- Add call to get the file and get the latest version
- Add call to download the Jar from the server
  - AKA from here: `https://maven.fabricmc.net/net/fabricmc/fabric-installer/0.11.1/fabric-installer-0.11.1.jar`
- Add checker to make sure java is installed.
  - Give user message on how to install Java.
- Add install call of the jar:
  - AKA command: `java -jar './fabric-installer-0.11.1.jar' client`
