package %%modgroup%%.fabric;

import %%modgroup%%.fabriclike.ExampleModFabricLike;
import net.fabricmc.api.ModInitializer;

public class %%modclass%%Fabric implements ModInitializer {
    @Override
    public void onInitialize() {
        %%modclass%%FabricLike.init();
    }
}
