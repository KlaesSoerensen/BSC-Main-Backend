namespace BSC_Main_Backend.dto.response;

public record PlayerInfoResponseDTO(uint Id, string IGN, uint sprite, List<uint> Achievements, bool HasCompletedTutorial);