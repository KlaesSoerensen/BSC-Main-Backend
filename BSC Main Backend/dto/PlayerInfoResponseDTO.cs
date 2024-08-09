namespace BSC_Main_Backend.dto;

public record PlayerInfoResponseDTO(uint Id, string IGN, uint sprite, List<uint> Achievements, bool HasCompletedTutorial);