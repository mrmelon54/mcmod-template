package %%modgroup%%.forge;

import %%modgroup%%.%%modclass%%;
import com.mrmelon54.omniver.forge.ConfigGui;
import dev.architectury.platform.forge.EventBuses;
import net.minecraftforge.fml.common.Mod;
import net.minecraftforge.fml.javafmlmod.FMLJavaModLoadingContext;

@Mod(%%modclass%%.MOD_ID)
public class %%modclass%%Forge {
    public %%modclass%%Forge() {
        EventBuses.registerModEventBus(%%modclass%%.MOD_ID, FMLJavaModLoadingContext.get().getModEventBus());
        ConfigGui.register((mc, screen) -> %%modclass%%.createConfigScreen(screen).get());
        %%modclass%%.init();
    }
}
