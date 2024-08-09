namespace BSC_Main_Backend.dto.response;

public record MiniGameInfoResponseDTO(
        uint Id,
        string Name,
        uint Icon,
        string Description,
        string Difficualties // (json)
    );