using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto;

public record GetAssetLODResponseDTO(
        uint Id,
        uint DetailLevel,
        Blob Blob // Maybe uint8[]
    );