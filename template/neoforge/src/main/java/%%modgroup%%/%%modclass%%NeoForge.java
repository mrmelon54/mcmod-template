package %%modgroup%%.neoforge;

import %%modgroup%%.%%modclass%%;
import net.neoforged.fml.ModLoadingContext;
import net.neoforged.fml.common.Mod;
import net.neoforged.fml.javafmlmod.FMLJavaModLoadingContext;
import net.neoforged.neoforge.client.ConfigScreenHandler.ConfigScreenFactory;

@Mod(%%modclass%%.MOD_ID)
public class %%modclass%%NeoForge {
    public %%modclass%%NeoForge() {
        ModLoadingContext.get().registerExtensionPoint(ConfigScreenFactory.class, () -> new ConfigScreenFactory((mc, screen) -> %%modclass%%.createConfigScreen(screen).get()));
        %%modclass%%.init();
    }
}
