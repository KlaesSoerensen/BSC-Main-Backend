namespace BSC_Main_Backend.dto.response;
using System.Collections.Generic;

public record ColonyInfoResponseDTO(
        uint Id,
        uint AccLevel,
        List<uint> Assets,
        List<uint> Locations
        );