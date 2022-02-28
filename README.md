# ForgeCLI

Package was created with the express intent to remove the guess work out of Mod acquiring and updating.

## Getting your Forge API Key

This is more complicated because you will be pulling/using the latest mod for the release of your game. To get started make sure you have a [CursedForge API Key](https://docs.curseforge.com/#getting-started). Then use it as a parameter for your build

## Quick start

Simple command to download latest fabric modules:

```bash
. ./forgecli.exe -forgekey '$2a$10...' -projects "416089,391366,552655" -family "fabric" -debug
```

## ForgeCLI Usage Details

**NOTE:** Due to the lack of Version pinning this can lead to unexpected behavior if the publisher updates mod unexpectedly.

To get started make sure you have a [CursedForge API Key](https://docs.curseforge.com/#getting-started).

**Default Folder Locations**:

- `Windows`: - "%AppData%/Roaming/.minecraft/mods"
- `Mac`: - "~/Library/Application Support/minecraft/mods"
- `Linux`: - "~/Library/Application Support/minecraft/mods"

**Mod basics**

- `Mod Release types`: - Release, Beta, and Alpha.
- `Mod dependencies`: - Required and Optional
- `Mod family`: - Fabric, Forge, and Bukkit.
- `Mod MC Version`: - 1.12.2, 1.18.2, etc.

**CLI Parameters**:

- `forgekey`: - Required as a parameter, or as an ENV var of FORGEKEY and MODS_FORGEAPI_KEY
- `file`: - Specially formatted json to manager larger sets of mods.
- `projects`: - Project ids that can be easily obtained from the Forge itself.
- `destination`: - Default is OS Client Mod folder, Target folder for the downloaded mods.
- `family`: - Used to filter mods based on server type. Options are Forge, Fabric, and Bukkit
- `release`: - Default is Release, Used to allow for Beta and Alpha mod downloads.
- `version`: - Default is LATEST, but this is Minecraft VERSION. e.g. 1.18.2
- `clearMods`: - Default is false, allows CLI to remove all mods before downloading new Mods.
- `downloadDependencies`: - Default is True, this uses the mods required dependencies to download missing mods.
- `debug`: - Enable extra logging.

## JSON File usage

**Field Description**:

- `name`: - is currently unused, but can be used to document each entry.
- `projectID`: - is the id found on the CurseForge website for a particular mod
- `releaseType`: - Type corresponds to forge's R, B, A icon for each file. Default Release, options are (release|beta|alpha).
- `fileName`: - is used for version pinning if latest file will not work for you.

**Example json File Format**:

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
    "releaseType": "beta"
  },
  {
    "name": "Biomes o plenty",
    "projectID": "220318",
    "fileName": "BiomesOPlenty-1.18.1-15.0.0.100-universal.jar",
    "releaseType": "release"
  }
]
```

### Manually Building and Testing

Make a `./.env` file in the root folder and add your forge key.

```text
FORGEKEY='$2a$10...'
```

CMD To test the build

```bash
go test ./...
```

### TODO List

- Update and proof read documentation
- Add normal info logging. Currently the app runs silent without -debug
- add fileName filter from json mods. Currently we have no method to pin versions
