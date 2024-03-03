package %%modgroup%%;

import me.shedaniel.autoconfig.ConfigData;
import me.shedaniel.autoconfig.annotation.Config;
import me.shedaniel.autoconfig.annotation.ConfigEntry;

@Config(name = "%%moddash%%")
@Config.Gui.Background("minecraft:textures/block/gold_block.png")
public class ConfigStructure implements ConfigData {
    public boolean modeEnabled = true;
}
