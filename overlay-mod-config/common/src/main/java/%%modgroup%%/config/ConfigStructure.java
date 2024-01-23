package %%modgroup%%.config;

import me.shedaniel.autoconfig.ConfigData;
import me.shedaniel.autoconfig.annotation.Config;
import me.shedaniel.autoconfig.annotation.ConfigEntry;

@Config(name = "%%modid%%")
@Config.Gui.Background("minecraft:textures/block/oak_planks.png")
public class ConfigStructure implements ConfigData {
    boolean modeEnabled = false;
}
