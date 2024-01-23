package %%modgroup%%.forge;

import dev.architectury.platform.forge.EventBuses;
import %%modgroup%%.%%modclass%%;
import net.minecraftforge.fml.common.Mod;
import net.minecraftforge.fml.javafmlmod.FMLJavaModLoadingContext;

@Mod(%%modclass%%.MOD_ID)
public class %%modclass%%Forge {
    public %%modclass%%Forge() {
        // Submit our event bus to let architectury register our content on the right time
        EventBuses.registerModEventBus(%%modclass%%.MOD_ID, FMLJavaModLoadingContext.get().getModEventBus());
        ModLoadingContext.get().registerExtensionPoint(ConfigScreenFactory.class, () -> new ConfigScreenFactory((mc, screen) -> %%modclass%%.createConfigScreen(screen)));
        %%modclass%%.init();
    }
}
