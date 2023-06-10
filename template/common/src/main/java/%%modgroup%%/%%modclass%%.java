package %%modgroup%%;

import dev.architectury.registry.registries.reg;

public class %%modclass%% {
    public static final String MOD_ID = "%%modid%%";

    public static void init() {
        System.out.println(ExampleExpectPlatform.getConfigDirectory().toAbsolutePath().normalize().toString());
    }
}
