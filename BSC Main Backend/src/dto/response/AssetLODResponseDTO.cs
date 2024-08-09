using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto;

public record AssetLODResponseDTO(
        uint Id,
        uint DetailLevel,
        Blob Blob // Maybe uint8[]
    );