using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto.response;

public record AssetLODResponseDTO(
        uint Id,
        uint DetailLevel,
        Blob Blob // Maybe uint8[]
    );