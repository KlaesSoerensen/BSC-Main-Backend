namespace BSC_Main_Backend.dto;
using System.Collections.Generic;

public record ColonyInfoResponseDTO(
        uint Id,
        uint AccLevel,
        List<uint> Assets,
        List<uint> Locations
        );