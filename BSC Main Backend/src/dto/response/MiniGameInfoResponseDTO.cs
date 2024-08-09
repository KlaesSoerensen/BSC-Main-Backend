namespace BSC_Main_Backend.dto;

public record MiniGameInfoResponseDTO(
        uint Id,
        string Name,
        uint Icon,
        string Description,
        string Difficualties // (json)
    );