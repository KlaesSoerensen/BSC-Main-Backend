using System.Reflection.Metadata;

namespace BSC_Main_Backend.dto;

public record FetchAssetResponseDTO(
        uint Id,
        string UseCase, // useCase: string "icon", "environment", "player", and so on.
        string Type, // Type: string"image/png", "image/jpeg", "image/webp", "image/svg+xml" (MIME types)
        uint Width, // res. x of LOD 0
        uint Height, // res. y of LOD 0 
        bool HasLODs, 
        Blob Blob, // Maybe uint8[]
        List<GetAssetLODResponseDTO> LODs
        );