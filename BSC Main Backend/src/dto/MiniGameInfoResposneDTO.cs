namespace BSC_Main_Backend.dto;

public record MiniGameInfoResposneDTO(
        uint Id,
        string Name,
        uint Icon,
        string Description,
        string Difficualties // (json)
    );