#!/usr/bin/python3
import tarfile
from pathlib import Path
import os.path
import urllib.request
import shutil
import sys
import json

def ask_for_info(prompt):
  while True:
    v = input(prompt)
    if v != "":
      return v

dl_url = "https://copy.mrmelon54.com/assets/arch-1.20.tar.gz"

modinfo = {
  "_example": "Delete this line when done",
  "modid": "clock-hud",
  "modname": "Clock HUD",
  "modclass": "ClockHud",
  "modgroup": "com.mrmelon54.ClockHud",
  "moddesc": "Enter the mod description",
  "modwebsite": "https://mrmelon54.com/minecraft/clock-hud",
  "modsource": "https://github.com/MrMelon54/clock-hud",
  "modissue": "https://github.com/MrMelon54/clock-hud/issues"
}

if not os.path.exists('mod-info.json'):
  with open('mod-info.json', 'w', encoding='utf8') as f:
    json.dump(modinfo, f, indent=2)
  print("Please fill out 'mod-info.json'")
  sys.exit(1)

with open('mod-info.json', 'r', encoding='utf8') as f:
  modinfo = json.load(f)

if "_example" in modinfo:
  print("Please fill out 'mod-info.json' and remove the '_example' field")
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

with urllib.request.urlopen(dl_url) as res, open("files.tar.gz", 'wb') as out:
  shutil.copyfileobj(res, out)

with tarfile.open('files.tar.gz', 'r') as tf:
  for member in tf.getmembers():
    if not member.isdir():
      mf = replace_mod_info_in_path(member.name)
      md = os.path.dirname(mf)
      print("Writing file:", mf)
      Path(md).mkdir(parents=True, exist_ok=True)
      f = tf.extractfile(member)
      if member.name.endswith('.jar'):
        with open(mf, "wb") as f2:
          f2.write(f.read())
      else:
        with open(mf, "w", encoding='utf8') as f2:
          c = f.read().decode('utf-8')
          f2.write(replace_mod_info_in_file(c))
