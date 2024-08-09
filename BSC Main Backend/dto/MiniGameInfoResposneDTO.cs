namespace BSC_Main_Backend.dto;

public record MiniGameInfoResposneDTO(
        uint Id,
        string Name,
        uint Icon,
        string Description,
        Object Difficualties //TODO: What Type here
    );