package %%modgroup%%.quilt;

import %%modgroup%%.fabriclike.%%modclass%%FabricLike;
import org.quiltmc.loader.api.ModContainer;
import org.quiltmc.qsl.base.api.entrypoint.ModInitializer;

public class %%modclass%%Quilt implements ModInitializer {
    @Override
    public void onInitialize(ModContainer mod) {
        %%modclass%%FabricLike.init();
    }
}
