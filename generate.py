#!/usr/bin/env python3
import tarfile
from pathlib import Path
import os.path
import sys
import json
import os
import untangle

if __name__ != "__main__":
    print("This is not a library")
    sys.exit(1)

script_path = os.path.realpath(__file__)
files_source = os.path.join(os.path.dirname(script_path), "files.tar.gz")


def ask_for_info(prompt):
    while True:
        v = input(prompt)
        if v != "":
            return v


def getLatestArchitecturyPlugin():
    url = "https://maven.architectury.dev/architectury-plugin/architectury-plugin.gradle.plugin/maven-metadata.xml"
    obj = untangle.parse(url)
    return obj.metadata.versioning.latest.cdata


def getLatestArchitecturyLoom():
    url = "https://maven.architectury.dev/dev/architectury/architectury-loom/maven-metadata.xml"
    obj = untangle.parse(url)
    return obj.metadata.versioning.latest.cdata


modinfo = {
    "_example": "Delete this line when done",
    "modid": "mod_id",
    "modname": "Mod Name",
    "modclass": "ModClass",
    "modgroup": "com.example.ModGroup",
    "moddesc": "Enter the mod description",
    "modwebsite": "https://mrmelon54.com/minecraft/mod-id",
    "modsource": "https://github.com/mrmelon54/mod_id",
    "modissue": "https://github.com/mrmelon54/mod_id/issues",
}

if not os.path.exists("mod-info.json"):
    with open("mod-info.json", "w", encoding="utf8") as f:
        json.dump(modinfo, f, indent=2)
    print("Please fill out 'mod-info.json'")
    sys.exit(1)

with open("mod-info.json", "r", encoding="utf8") as f:
    modinfo = json.load(f)

modinfo["moddash"] = modinfo["modid"]
modinfo["minecraft_version"] = ask_for_info("Minecraft version: ")
modinfo["architectury_plugin_version"] = getLatestArchitecturyPlugin()
modinfo["architectury_loom_version"] = getLatestArchitecturyLoom()

if "_example" in modinfo:
    print("Please fill out 'mod-info.json' and remove the '_example' field")
    sys.exit(1)

if "-" in modinfo["modid"]:
    print("Forge doesn't allow '-' so use '_' instead")
    sys.exit(1)


def replace_mod_info_in_path(x):
    for k in modinfo:
        if k == "modgroup":
            x = x.replace("%%" + k + "%%", modinfo[k].replace(".", "/"))
        else:
            x = x.replace("%%" + k + "%%", modinfo[k])
    return x


def replace_mod_info_in_file(x):
    for k in modinfo:
        x = x.replace("%%" + k + "%%", modinfo[k])
    return x


overrideFiles = False
with tarfile.open(files_source, "r") as tf:
    for member in tf.getmembers():
        if not member.isdir():
            mf = replace_mod_info_in_path(member.name)
            md = os.path.dirname(mf)
            print("Writing file:", mf)
            Path(md).mkdir(parents=True, exist_ok=True)
            f = tf.extractfile(member)
            if os.path.exists(mf) and not overrideFiles:
                con = input("Trying to override file, ignore existing files? [y/N]: ")
                if con.lower() in ["y", "yes"]:
                    overrideFiles = True
                else:
                    print("Goodbye")
                    sys.exit(1)
            if member.name.endswith(".jar"):
                with open(mf, "wb") as f2:
                    f2.write(f.read())
            else:
                with open(mf, "w", encoding="utf8") as f2:
                    c = f.read().decode("utf-8")
                    f2.write(replace_mod_info_in_file(c))

os.system("chmod +x gradlew")
