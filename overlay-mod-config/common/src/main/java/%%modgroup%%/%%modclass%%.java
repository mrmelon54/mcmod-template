package %%modgroup%%;

public class %%modclass%% {
    public static final String MOD_ID = "%%modid%%";
    public static ConfigStructure CONFIG = AutoConfig.register(ConfigStructure.class).get();

    public static void init() {
    }

    public static Screen createConfigScreen(Screen parent) {
        return AutoConfig.getConfigScreen(ConfigStructure.class, parent).get();
    }
}
